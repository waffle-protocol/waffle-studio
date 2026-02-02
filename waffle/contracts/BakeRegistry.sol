// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title BakeRegistry
 * @dev Manages Bake requests with SYRUP token escrow for the Waffle P2P marketplace
 * @notice Users create requests with SYRUP bounty, Bakers submit solutions, users accept/reject
 */
contract BakeRegistry is ReentrancyGuard {
    using SafeERC20 for IERC20;

    // ============================================================================
    // Types
    // ============================================================================

    enum RequestStatus {
        PENDING,    // Request created, waiting for Baker
        SUBMITTED,  // Baker submitted solution
        ACCEPTED,   // User accepted solution, reward paid
        REJECTED,   // User rejected solution
        CANCELLED   // Request cancelled by requester
    }

    struct BakeRequest {
        address requester;      // User who created request
        address baker;          // Baker who submitted solution
        bytes32 codeHash;       // Hash of code being baked
        bytes32 solutionHash;   // Hash of solution (if submitted)
        uint256 reward;         // SYRUP amount escrowed
        uint256 tokenUsage;     // Usage reported by provider
        uint256 createdAt;      // Timestamp of creation
        RequestStatus status;   // Current status
    }

    // ============================================================================
    // State Variables
    // ============================================================================

    IERC20 public immutable syrupToken;
    
    uint256 public nextRequestId;
    mapping(uint256 => BakeRequest) public requests;
    
    // Track requests by user
    mapping(address => uint256[]) public userRequests;
    // Track requests by baker
    mapping(address => uint256[]) public bakerSubmissions;

    // ============================================================================
    // Events
    // ============================================================================

    event RequestCreated(
        uint256 indexed requestId,
        address indexed requester,
        bytes32 codeHash,
        uint256 reward
    );

    event SolutionSubmitted(
        uint256 indexed requestId,
        address indexed baker,
        bytes32 solutionHash,
        uint256 tokenUsage
    );

    event SolutionAccepted(
        uint256 indexed requestId,
        address indexed baker,
        uint256 rewardPaid,
        uint256 refundAmount
    );

    event SolutionRejected(
        uint256 indexed requestId,
        address indexed baker
    );

    event RequestCancelled(
        uint256 indexed requestId,
        address indexed requester,
        uint256 refundAmount
    );

    // ============================================================================
    // Errors
    // ============================================================================

    error InvalidReward();
    error RequestNotFound();
    error NotRequester();
    error NotPending();
    error NotSubmitted();
    error AlreadySubmitted();
    error CannotSubmitOwnRequest();
    error InvalidPaymentAmount();

    // ============================================================================
    // Constructor
    // ============================================================================

    /**
     * @dev Sets the SYRUP token address
     * @param _syrupToken Address of the SyrupToken contract
     */
    constructor(address _syrupToken) {
        syrupToken = IERC20(_syrupToken);
    }

    // ============================================================================
    // External Functions
    // ============================================================================

    /**
     * @dev Create a new Bake request with SYRUP bounty
     * @param codeHash Hash of the code to be baked
     * @param reward Amount of SYRUP tokens to escrow as reward
     * @return requestId The ID of the created request
     */
    function createRequest(bytes32 codeHash, uint256 reward) 
        external 
        nonReentrant 
        returns (uint256 requestId) 
    {
        if (reward == 0) revert InvalidReward();

        // Transfer SYRUP from user to this contract (escrow)
        syrupToken.safeTransferFrom(msg.sender, address(this), reward);

        requestId = nextRequestId++;
        
        requests[requestId] = BakeRequest({
            requester: msg.sender,
            baker: address(0),
            codeHash: codeHash,
            solutionHash: bytes32(0),
            reward: reward,
            tokenUsage: 0,
            createdAt: block.timestamp,
            status: RequestStatus.PENDING
        });

        userRequests[msg.sender].push(requestId);

        emit RequestCreated(requestId, msg.sender, codeHash, reward);
    }

    /**
     * @dev Baker submits a solution for a pending request
     * @param requestId ID of the request
     * @param solutionHash Hash of the solution
     * @param tokenUsage Token usage reported by provider
     */
    function submitSolution(uint256 requestId, bytes32 solutionHash, uint256 tokenUsage) 
        external 
    {
        BakeRequest storage request = requests[requestId];
        
        if (request.requester == address(0)) revert RequestNotFound();
        if (request.status != RequestStatus.PENDING) revert NotPending();
        if (request.requester == msg.sender) revert CannotSubmitOwnRequest();

        request.baker = msg.sender;
        request.solutionHash = solutionHash;
        request.tokenUsage = tokenUsage;
        request.status = RequestStatus.SUBMITTED;

        bakerSubmissions[msg.sender].push(requestId);

        emit SolutionSubmitted(requestId, msg.sender, solutionHash, tokenUsage);
    }

    /**
     * @dev Requester accepts the solution and pays the Baker
     * @param requestId ID of the request
     * @param paymentAmount Amount to pay the baker (must be <= escrowed reward)
     */
    function acceptSolution(uint256 requestId, uint256 paymentAmount) 
        external 
        nonReentrant 
    {
        BakeRequest storage request = requests[requestId];
        
        if (request.requester != msg.sender) revert NotRequester();
        if (request.status != RequestStatus.SUBMITTED) revert NotSubmitted();
        if (paymentAmount > request.reward) revert InvalidPaymentAmount();

        request.status = RequestStatus.ACCEPTED;

        // Transfer payment to Baker
        if (paymentAmount > 0) {
            syrupToken.safeTransfer(request.baker, paymentAmount);
        }

        // Refund remaining to Requester
        uint256 refund = request.reward - paymentAmount;
        if (refund > 0) {
            syrupToken.safeTransfer(request.requester, refund);
        }

        emit SolutionAccepted(requestId, request.baker, paymentAmount, refund);
    }

    /**
     * @dev Requester rejects the solution, request returns to PENDING
     * @param requestId ID of the request
     */
    function rejectSolution(uint256 requestId) 
        external 
    {
        BakeRequest storage request = requests[requestId];
        
        if (request.requester != msg.sender) revert NotRequester();
        if (request.status != RequestStatus.SUBMITTED) revert NotSubmitted();

        // Reset to pending for another baker to try
        request.status = RequestStatus.REJECTED;

        emit SolutionRejected(requestId, request.baker);
    }

    /**
     * @dev Requester cancels a pending request and gets refund
     * @param requestId ID of the request
     */
    function cancelRequest(uint256 requestId) 
        external 
        nonReentrant 
    {
        BakeRequest storage request = requests[requestId];
        
        if (request.requester != msg.sender) revert NotRequester();
        if (request.status != RequestStatus.PENDING) revert NotPending();

        request.status = RequestStatus.CANCELLED;

        // Refund SYRUP to requester
        syrupToken.safeTransfer(msg.sender, request.reward);

        emit RequestCancelled(requestId, msg.sender, request.reward);
    }

    // ============================================================================
    // View Functions
    // ============================================================================

    /**
     * @dev Get request details
     */
    function getRequest(uint256 requestId) 
        external 
        view 
        returns (BakeRequest memory) 
    {
        return requests[requestId];
    }

    /**
     * @dev Get all request IDs for a user
     */
    function getUserRequests(address user) 
        external 
        view 
        returns (uint256[] memory) 
    {
        return userRequests[user];
    }

    /**
     * @dev Get all submission IDs for a baker
     */
    function getBakerSubmissions(address baker) 
        external 
        view 
        returns (uint256[] memory) 
    {
        return bakerSubmissions[baker];
    }

    /**
     * @dev Get total number of requests
     */
    function totalRequests() external view returns (uint256) {
        return nextRequestId;
    }
}
