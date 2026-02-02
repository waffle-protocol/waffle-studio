// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BakeRegistryBakeRequest is an auto generated low-level Go binding around an user-defined struct.
type BakeRegistryBakeRequest struct {
	Requester    common.Address
	Baker        common.Address
	CodeHash     [32]byte
	SolutionHash [32]byte
	Reward       *big.Int
	TokenUsage   *big.Int
	CreatedAt    *big.Int
	Status       uint8
}

// BakeRegistryMetaData contains all meta data concerning the BakeRegistry contract.
var BakeRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_syrupToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptSolution\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"paymentAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bakerSubmissions\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelRequest\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createRequest\",\"inputs\":[{\"name\":\"codeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"reward\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getBakerSubmissions\",\"inputs\":[{\"name\":\"baker\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRequest\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structBakeRegistry.BakeRequest\",\"components\":[{\"name\":\"requester\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"baker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"codeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"solutionHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"reward\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenUsage\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumBakeRegistry.RequestStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUserRequests\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextRequestId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rejectSolution\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requests\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"requester\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"baker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"codeHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"solutionHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"reward\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"tokenUsage\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumBakeRegistry.RequestStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"submitSolution\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"solutionHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"tokenUsage\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"syrupToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalRequests\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userRequests\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"RequestCancelled\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"refundAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestCreated\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"requester\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"codeHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"reward\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SolutionAccepted\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"baker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"rewardPaid\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"refundAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SolutionRejected\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"baker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SolutionSubmitted\",\"inputs\":[{\"name\":\"requestId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"baker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"solutionHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"tokenUsage\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadySubmitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CannotSubmitOwnRequest\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPaymentAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReward\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotPending\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRequester\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSubmitted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RequestNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// BakeRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use BakeRegistryMetaData.ABI instead.
var BakeRegistryABI = BakeRegistryMetaData.ABI

// BakeRegistry is an auto generated Go binding around an Ethereum contract.
type BakeRegistry struct {
	BakeRegistryCaller     // Read-only binding to the contract
	BakeRegistryTransactor // Write-only binding to the contract
	BakeRegistryFilterer   // Log filterer for contract events
}

// BakeRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type BakeRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BakeRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BakeRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BakeRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BakeRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BakeRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BakeRegistrySession struct {
	Contract     *BakeRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BakeRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BakeRegistryCallerSession struct {
	Contract *BakeRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BakeRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BakeRegistryTransactorSession struct {
	Contract     *BakeRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BakeRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type BakeRegistryRaw struct {
	Contract *BakeRegistry // Generic contract binding to access the raw methods on
}

// BakeRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BakeRegistryCallerRaw struct {
	Contract *BakeRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// BakeRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BakeRegistryTransactorRaw struct {
	Contract *BakeRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBakeRegistry creates a new instance of BakeRegistry, bound to a specific deployed contract.
func NewBakeRegistry(address common.Address, backend bind.ContractBackend) (*BakeRegistry, error) {
	contract, err := bindBakeRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BakeRegistry{BakeRegistryCaller: BakeRegistryCaller{contract: contract}, BakeRegistryTransactor: BakeRegistryTransactor{contract: contract}, BakeRegistryFilterer: BakeRegistryFilterer{contract: contract}}, nil
}

// NewBakeRegistryCaller creates a new read-only instance of BakeRegistry, bound to a specific deployed contract.
func NewBakeRegistryCaller(address common.Address, caller bind.ContractCaller) (*BakeRegistryCaller, error) {
	contract, err := bindBakeRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BakeRegistryCaller{contract: contract}, nil
}

// NewBakeRegistryTransactor creates a new write-only instance of BakeRegistry, bound to a specific deployed contract.
func NewBakeRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*BakeRegistryTransactor, error) {
	contract, err := bindBakeRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BakeRegistryTransactor{contract: contract}, nil
}

// NewBakeRegistryFilterer creates a new log filterer instance of BakeRegistry, bound to a specific deployed contract.
func NewBakeRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*BakeRegistryFilterer, error) {
	contract, err := bindBakeRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BakeRegistryFilterer{contract: contract}, nil
}

// bindBakeRegistry binds a generic wrapper to an already deployed contract.
func bindBakeRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BakeRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BakeRegistry *BakeRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BakeRegistry.Contract.BakeRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BakeRegistry *BakeRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BakeRegistry.Contract.BakeRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BakeRegistry *BakeRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BakeRegistry.Contract.BakeRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BakeRegistry *BakeRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BakeRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BakeRegistry *BakeRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BakeRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BakeRegistry *BakeRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BakeRegistry.Contract.contract.Transact(opts, method, params...)
}

// BakerSubmissions is a free data retrieval call binding the contract method 0x095afe1d.
//
// Solidity: function bakerSubmissions(address , uint256 ) view returns(uint256)
func (_BakeRegistry *BakeRegistryCaller) BakerSubmissions(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BakeRegistry.contract.Call(opts, &out, "bakerSubmissions", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BakerSubmissions is a free data retrieval call binding the contract method 0x095afe1d.
//
// Solidity: function bakerSubmissions(address , uint256 ) view returns(uint256)
func (_BakeRegistry *BakeRegistrySession) BakerSubmissions(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _BakeRegistry.Contract.BakerSubmissions(&_BakeRegistry.CallOpts, arg0, arg1)
}

// BakerSubmissions is a free data retrieval call binding the contract method 0x095afe1d.
//
// Solidity: function bakerSubmissions(address , uint256 ) view returns(uint256)
func (_BakeRegistry *BakeRegistryCallerSession) BakerSubmissions(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _BakeRegistry.Contract.BakerSubmissions(&_BakeRegistry.CallOpts, arg0, arg1)
}

// GetBakerSubmissions is a free data retrieval call binding the contract method 0x8868c48c.
//
// Solidity: function getBakerSubmissions(address baker) view returns(uint256[])
func (_BakeRegistry *BakeRegistryCaller) GetBakerSubmissions(opts *bind.CallOpts, baker common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _BakeRegistry.contract.Call(opts, &out, "getBakerSubmissions", baker)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetBakerSubmissions is a free data retrieval call binding the contract method 0x8868c48c.
//
// Solidity: function getBakerSubmissions(address baker) view returns(uint256[])
func (_BakeRegistry *BakeRegistrySession) GetBakerSubmissions(baker common.Address) ([]*big.Int, error) {
	return _BakeRegistry.Contract.GetBakerSubmissions(&_BakeRegistry.CallOpts, baker)
}

// GetBakerSubmissions is a free data retrieval call binding the contract method 0x8868c48c.
//
// Solidity: function getBakerSubmissions(address baker) view returns(uint256[])
func (_BakeRegistry *BakeRegistryCallerSession) GetBakerSubmissions(baker common.Address) ([]*big.Int, error) {
	return _BakeRegistry.Contract.GetBakerSubmissions(&_BakeRegistry.CallOpts, baker)
}

// GetRequest is a free data retrieval call binding the contract method 0xc58343ef.
//
// Solidity: function getRequest(uint256 requestId) view returns((address,address,bytes32,bytes32,uint256,uint256,uint256,uint8))
func (_BakeRegistry *BakeRegistryCaller) GetRequest(opts *bind.CallOpts, requestId *big.Int) (BakeRegistryBakeRequest, error) {
	var out []interface{}
	err := _BakeRegistry.contract.Call(opts, &out, "getRequest", requestId)

	if err != nil {
		return *new(BakeRegistryBakeRequest), err
	}

	out0 := *abi.ConvertType(out[0], new(BakeRegistryBakeRequest)).(*BakeRegistryBakeRequest)

	return out0, err

}

// GetRequest is a free data retrieval call binding the contract method 0xc58343ef.
//
// Solidity: function getRequest(uint256 requestId) view returns((address,address,bytes32,bytes32,uint256,uint256,uint256,uint8))
func (_BakeRegistry *BakeRegistrySession) GetRequest(requestId *big.Int) (BakeRegistryBakeRequest, error) {
	return _BakeRegistry.Contract.GetRequest(&_BakeRegistry.CallOpts, requestId)
}

// GetRequest is a free data retrieval call binding the contract method 0xc58343ef.
//
// Solidity: function getRequest(uint256 requestId) view returns((address,address,bytes32,bytes32,uint256,uint256,uint256,uint8))
func (_BakeRegistry *BakeRegistryCallerSession) GetRequest(requestId *big.Int) (BakeRegistryBakeRequest, error) {
	return _BakeRegistry.Contract.GetRequest(&_BakeRegistry.CallOpts, requestId)
}

// GetUserRequests is a free data retrieval call binding the contract method 0xb337cf74.
//
// Solidity: function getUserRequests(address user) view returns(uint256[])
func (_BakeRegistry *BakeRegistryCaller) GetUserRequests(opts *bind.CallOpts, user common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _BakeRegistry.contract.Call(opts, &out, "getUserRequests", user)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUserRequests is a free data retrieval call binding the contract method 0xb337cf74.
//
// Solidity: function getUserRequests(address user) view returns(uint256[])
func (_BakeRegistry *BakeRegistrySession) GetUserRequests(user common.Address) ([]*big.Int, error) {
	return _BakeRegistry.Contract.GetUserRequests(&_BakeRegistry.CallOpts, user)
}

// GetUserRequests is a free data retrieval call binding the contract method 0xb337cf74.
//
// Solidity: function getUserRequests(address user) view returns(uint256[])
func (_BakeRegistry *BakeRegistryCallerSession) GetUserRequests(user common.Address) ([]*big.Int, error) {
	return _BakeRegistry.Contract.GetUserRequests(&_BakeRegistry.CallOpts, user)
}

// NextRequestId is a free data retrieval call binding the contract method 0x6a84a985.
//
// Solidity: function nextRequestId() view returns(uint256)
func (_BakeRegistry *BakeRegistryCaller) NextRequestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BakeRegistry.contract.Call(opts, &out, "nextRequestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextRequestId is a free data retrieval call binding the contract method 0x6a84a985.
//
// Solidity: function nextRequestId() view returns(uint256)
func (_BakeRegistry *BakeRegistrySession) NextRequestId() (*big.Int, error) {
	return _BakeRegistry.Contract.NextRequestId(&_BakeRegistry.CallOpts)
}

// NextRequestId is a free data retrieval call binding the contract method 0x6a84a985.
//
// Solidity: function nextRequestId() view returns(uint256)
func (_BakeRegistry *BakeRegistryCallerSession) NextRequestId() (*big.Int, error) {
	return _BakeRegistry.Contract.NextRequestId(&_BakeRegistry.CallOpts)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(address requester, address baker, bytes32 codeHash, bytes32 solutionHash, uint256 reward, uint256 tokenUsage, uint256 createdAt, uint8 status)
func (_BakeRegistry *BakeRegistryCaller) Requests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Requester    common.Address
	Baker        common.Address
	CodeHash     [32]byte
	SolutionHash [32]byte
	Reward       *big.Int
	TokenUsage   *big.Int
	CreatedAt    *big.Int
	Status       uint8
}, error) {
	var out []interface{}
	err := _BakeRegistry.contract.Call(opts, &out, "requests", arg0)

	outstruct := new(struct {
		Requester    common.Address
		Baker        common.Address
		CodeHash     [32]byte
		SolutionHash [32]byte
		Reward       *big.Int
		TokenUsage   *big.Int
		CreatedAt    *big.Int
		Status       uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Requester = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Baker = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.CodeHash = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.SolutionHash = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Reward = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.TokenUsage = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.CreatedAt = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[7], new(uint8)).(*uint8)

	return *outstruct, err

}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(address requester, address baker, bytes32 codeHash, bytes32 solutionHash, uint256 reward, uint256 tokenUsage, uint256 createdAt, uint8 status)
func (_BakeRegistry *BakeRegistrySession) Requests(arg0 *big.Int) (struct {
	Requester    common.Address
	Baker        common.Address
	CodeHash     [32]byte
	SolutionHash [32]byte
	Reward       *big.Int
	TokenUsage   *big.Int
	CreatedAt    *big.Int
	Status       uint8
}, error) {
	return _BakeRegistry.Contract.Requests(&_BakeRegistry.CallOpts, arg0)
}

// Requests is a free data retrieval call binding the contract method 0x81d12c58.
//
// Solidity: function requests(uint256 ) view returns(address requester, address baker, bytes32 codeHash, bytes32 solutionHash, uint256 reward, uint256 tokenUsage, uint256 createdAt, uint8 status)
func (_BakeRegistry *BakeRegistryCallerSession) Requests(arg0 *big.Int) (struct {
	Requester    common.Address
	Baker        common.Address
	CodeHash     [32]byte
	SolutionHash [32]byte
	Reward       *big.Int
	TokenUsage   *big.Int
	CreatedAt    *big.Int
	Status       uint8
}, error) {
	return _BakeRegistry.Contract.Requests(&_BakeRegistry.CallOpts, arg0)
}

// SyrupToken is a free data retrieval call binding the contract method 0xff7f3c51.
//
// Solidity: function syrupToken() view returns(address)
func (_BakeRegistry *BakeRegistryCaller) SyrupToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BakeRegistry.contract.Call(opts, &out, "syrupToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SyrupToken is a free data retrieval call binding the contract method 0xff7f3c51.
//
// Solidity: function syrupToken() view returns(address)
func (_BakeRegistry *BakeRegistrySession) SyrupToken() (common.Address, error) {
	return _BakeRegistry.Contract.SyrupToken(&_BakeRegistry.CallOpts)
}

// SyrupToken is a free data retrieval call binding the contract method 0xff7f3c51.
//
// Solidity: function syrupToken() view returns(address)
func (_BakeRegistry *BakeRegistryCallerSession) SyrupToken() (common.Address, error) {
	return _BakeRegistry.Contract.SyrupToken(&_BakeRegistry.CallOpts)
}

// TotalRequests is a free data retrieval call binding the contract method 0x8aea61dc.
//
// Solidity: function totalRequests() view returns(uint256)
func (_BakeRegistry *BakeRegistryCaller) TotalRequests(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BakeRegistry.contract.Call(opts, &out, "totalRequests")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalRequests is a free data retrieval call binding the contract method 0x8aea61dc.
//
// Solidity: function totalRequests() view returns(uint256)
func (_BakeRegistry *BakeRegistrySession) TotalRequests() (*big.Int, error) {
	return _BakeRegistry.Contract.TotalRequests(&_BakeRegistry.CallOpts)
}

// TotalRequests is a free data retrieval call binding the contract method 0x8aea61dc.
//
// Solidity: function totalRequests() view returns(uint256)
func (_BakeRegistry *BakeRegistryCallerSession) TotalRequests() (*big.Int, error) {
	return _BakeRegistry.Contract.TotalRequests(&_BakeRegistry.CallOpts)
}

// UserRequests is a free data retrieval call binding the contract method 0x263cb6b6.
//
// Solidity: function userRequests(address , uint256 ) view returns(uint256)
func (_BakeRegistry *BakeRegistryCaller) UserRequests(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BakeRegistry.contract.Call(opts, &out, "userRequests", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserRequests is a free data retrieval call binding the contract method 0x263cb6b6.
//
// Solidity: function userRequests(address , uint256 ) view returns(uint256)
func (_BakeRegistry *BakeRegistrySession) UserRequests(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _BakeRegistry.Contract.UserRequests(&_BakeRegistry.CallOpts, arg0, arg1)
}

// UserRequests is a free data retrieval call binding the contract method 0x263cb6b6.
//
// Solidity: function userRequests(address , uint256 ) view returns(uint256)
func (_BakeRegistry *BakeRegistryCallerSession) UserRequests(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _BakeRegistry.Contract.UserRequests(&_BakeRegistry.CallOpts, arg0, arg1)
}

// AcceptSolution is a paid mutator transaction binding the contract method 0x08bad193.
//
// Solidity: function acceptSolution(uint256 requestId, uint256 paymentAmount) returns()
func (_BakeRegistry *BakeRegistryTransactor) AcceptSolution(opts *bind.TransactOpts, requestId *big.Int, paymentAmount *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.contract.Transact(opts, "acceptSolution", requestId, paymentAmount)
}

// AcceptSolution is a paid mutator transaction binding the contract method 0x08bad193.
//
// Solidity: function acceptSolution(uint256 requestId, uint256 paymentAmount) returns()
func (_BakeRegistry *BakeRegistrySession) AcceptSolution(requestId *big.Int, paymentAmount *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.AcceptSolution(&_BakeRegistry.TransactOpts, requestId, paymentAmount)
}

// AcceptSolution is a paid mutator transaction binding the contract method 0x08bad193.
//
// Solidity: function acceptSolution(uint256 requestId, uint256 paymentAmount) returns()
func (_BakeRegistry *BakeRegistryTransactorSession) AcceptSolution(requestId *big.Int, paymentAmount *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.AcceptSolution(&_BakeRegistry.TransactOpts, requestId, paymentAmount)
}

// CancelRequest is a paid mutator transaction binding the contract method 0x3015394c.
//
// Solidity: function cancelRequest(uint256 requestId) returns()
func (_BakeRegistry *BakeRegistryTransactor) CancelRequest(opts *bind.TransactOpts, requestId *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.contract.Transact(opts, "cancelRequest", requestId)
}

// CancelRequest is a paid mutator transaction binding the contract method 0x3015394c.
//
// Solidity: function cancelRequest(uint256 requestId) returns()
func (_BakeRegistry *BakeRegistrySession) CancelRequest(requestId *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.CancelRequest(&_BakeRegistry.TransactOpts, requestId)
}

// CancelRequest is a paid mutator transaction binding the contract method 0x3015394c.
//
// Solidity: function cancelRequest(uint256 requestId) returns()
func (_BakeRegistry *BakeRegistryTransactorSession) CancelRequest(requestId *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.CancelRequest(&_BakeRegistry.TransactOpts, requestId)
}

// CreateRequest is a paid mutator transaction binding the contract method 0x3bd0c42e.
//
// Solidity: function createRequest(bytes32 codeHash, uint256 reward) returns(uint256 requestId)
func (_BakeRegistry *BakeRegistryTransactor) CreateRequest(opts *bind.TransactOpts, codeHash [32]byte, reward *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.contract.Transact(opts, "createRequest", codeHash, reward)
}

// CreateRequest is a paid mutator transaction binding the contract method 0x3bd0c42e.
//
// Solidity: function createRequest(bytes32 codeHash, uint256 reward) returns(uint256 requestId)
func (_BakeRegistry *BakeRegistrySession) CreateRequest(codeHash [32]byte, reward *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.CreateRequest(&_BakeRegistry.TransactOpts, codeHash, reward)
}

// CreateRequest is a paid mutator transaction binding the contract method 0x3bd0c42e.
//
// Solidity: function createRequest(bytes32 codeHash, uint256 reward) returns(uint256 requestId)
func (_BakeRegistry *BakeRegistryTransactorSession) CreateRequest(codeHash [32]byte, reward *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.CreateRequest(&_BakeRegistry.TransactOpts, codeHash, reward)
}

// RejectSolution is a paid mutator transaction binding the contract method 0xcda81874.
//
// Solidity: function rejectSolution(uint256 requestId) returns()
func (_BakeRegistry *BakeRegistryTransactor) RejectSolution(opts *bind.TransactOpts, requestId *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.contract.Transact(opts, "rejectSolution", requestId)
}

// RejectSolution is a paid mutator transaction binding the contract method 0xcda81874.
//
// Solidity: function rejectSolution(uint256 requestId) returns()
func (_BakeRegistry *BakeRegistrySession) RejectSolution(requestId *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.RejectSolution(&_BakeRegistry.TransactOpts, requestId)
}

// RejectSolution is a paid mutator transaction binding the contract method 0xcda81874.
//
// Solidity: function rejectSolution(uint256 requestId) returns()
func (_BakeRegistry *BakeRegistryTransactorSession) RejectSolution(requestId *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.RejectSolution(&_BakeRegistry.TransactOpts, requestId)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x0b8ddcb1.
//
// Solidity: function submitSolution(uint256 requestId, bytes32 solutionHash, uint256 tokenUsage) returns()
func (_BakeRegistry *BakeRegistryTransactor) SubmitSolution(opts *bind.TransactOpts, requestId *big.Int, solutionHash [32]byte, tokenUsage *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.contract.Transact(opts, "submitSolution", requestId, solutionHash, tokenUsage)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x0b8ddcb1.
//
// Solidity: function submitSolution(uint256 requestId, bytes32 solutionHash, uint256 tokenUsage) returns()
func (_BakeRegistry *BakeRegistrySession) SubmitSolution(requestId *big.Int, solutionHash [32]byte, tokenUsage *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.SubmitSolution(&_BakeRegistry.TransactOpts, requestId, solutionHash, tokenUsage)
}

// SubmitSolution is a paid mutator transaction binding the contract method 0x0b8ddcb1.
//
// Solidity: function submitSolution(uint256 requestId, bytes32 solutionHash, uint256 tokenUsage) returns()
func (_BakeRegistry *BakeRegistryTransactorSession) SubmitSolution(requestId *big.Int, solutionHash [32]byte, tokenUsage *big.Int) (*types.Transaction, error) {
	return _BakeRegistry.Contract.SubmitSolution(&_BakeRegistry.TransactOpts, requestId, solutionHash, tokenUsage)
}

// BakeRegistryRequestCancelledIterator is returned from FilterRequestCancelled and is used to iterate over the raw logs and unpacked data for RequestCancelled events raised by the BakeRegistry contract.
type BakeRegistryRequestCancelledIterator struct {
	Event *BakeRegistryRequestCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BakeRegistryRequestCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BakeRegistryRequestCancelled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BakeRegistryRequestCancelled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BakeRegistryRequestCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BakeRegistryRequestCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BakeRegistryRequestCancelled represents a RequestCancelled event raised by the BakeRegistry contract.
type BakeRegistryRequestCancelled struct {
	RequestId    *big.Int
	Requester    common.Address
	RefundAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRequestCancelled is a free log retrieval operation binding the contract event 0x3796cd2bbde8f185a2341d1c3fb8207b168febf5cd44d360661956cfdc70a177.
//
// Solidity: event RequestCancelled(uint256 indexed requestId, address indexed requester, uint256 refundAmount)
func (_BakeRegistry *BakeRegistryFilterer) FilterRequestCancelled(opts *bind.FilterOpts, requestId []*big.Int, requester []common.Address) (*BakeRegistryRequestCancelledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _BakeRegistry.contract.FilterLogs(opts, "RequestCancelled", requestIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &BakeRegistryRequestCancelledIterator{contract: _BakeRegistry.contract, event: "RequestCancelled", logs: logs, sub: sub}, nil
}

// WatchRequestCancelled is a free log subscription operation binding the contract event 0x3796cd2bbde8f185a2341d1c3fb8207b168febf5cd44d360661956cfdc70a177.
//
// Solidity: event RequestCancelled(uint256 indexed requestId, address indexed requester, uint256 refundAmount)
func (_BakeRegistry *BakeRegistryFilterer) WatchRequestCancelled(opts *bind.WatchOpts, sink chan<- *BakeRegistryRequestCancelled, requestId []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _BakeRegistry.contract.WatchLogs(opts, "RequestCancelled", requestIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BakeRegistryRequestCancelled)
				if err := _BakeRegistry.contract.UnpackLog(event, "RequestCancelled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestCancelled is a log parse operation binding the contract event 0x3796cd2bbde8f185a2341d1c3fb8207b168febf5cd44d360661956cfdc70a177.
//
// Solidity: event RequestCancelled(uint256 indexed requestId, address indexed requester, uint256 refundAmount)
func (_BakeRegistry *BakeRegistryFilterer) ParseRequestCancelled(log types.Log) (*BakeRegistryRequestCancelled, error) {
	event := new(BakeRegistryRequestCancelled)
	if err := _BakeRegistry.contract.UnpackLog(event, "RequestCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BakeRegistryRequestCreatedIterator is returned from FilterRequestCreated and is used to iterate over the raw logs and unpacked data for RequestCreated events raised by the BakeRegistry contract.
type BakeRegistryRequestCreatedIterator struct {
	Event *BakeRegistryRequestCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BakeRegistryRequestCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BakeRegistryRequestCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BakeRegistryRequestCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BakeRegistryRequestCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BakeRegistryRequestCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BakeRegistryRequestCreated represents a RequestCreated event raised by the BakeRegistry contract.
type BakeRegistryRequestCreated struct {
	RequestId *big.Int
	Requester common.Address
	CodeHash  [32]byte
	Reward    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestCreated is a free log retrieval operation binding the contract event 0x78e1c38f7bce169c7cf026c9115bab62243678331df819e47ba8f2cd48ba259b.
//
// Solidity: event RequestCreated(uint256 indexed requestId, address indexed requester, bytes32 codeHash, uint256 reward)
func (_BakeRegistry *BakeRegistryFilterer) FilterRequestCreated(opts *bind.FilterOpts, requestId []*big.Int, requester []common.Address) (*BakeRegistryRequestCreatedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _BakeRegistry.contract.FilterLogs(opts, "RequestCreated", requestIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return &BakeRegistryRequestCreatedIterator{contract: _BakeRegistry.contract, event: "RequestCreated", logs: logs, sub: sub}, nil
}

// WatchRequestCreated is a free log subscription operation binding the contract event 0x78e1c38f7bce169c7cf026c9115bab62243678331df819e47ba8f2cd48ba259b.
//
// Solidity: event RequestCreated(uint256 indexed requestId, address indexed requester, bytes32 codeHash, uint256 reward)
func (_BakeRegistry *BakeRegistryFilterer) WatchRequestCreated(opts *bind.WatchOpts, sink chan<- *BakeRegistryRequestCreated, requestId []*big.Int, requester []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var requesterRule []interface{}
	for _, requesterItem := range requester {
		requesterRule = append(requesterRule, requesterItem)
	}

	logs, sub, err := _BakeRegistry.contract.WatchLogs(opts, "RequestCreated", requestIdRule, requesterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BakeRegistryRequestCreated)
				if err := _BakeRegistry.contract.UnpackLog(event, "RequestCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestCreated is a log parse operation binding the contract event 0x78e1c38f7bce169c7cf026c9115bab62243678331df819e47ba8f2cd48ba259b.
//
// Solidity: event RequestCreated(uint256 indexed requestId, address indexed requester, bytes32 codeHash, uint256 reward)
func (_BakeRegistry *BakeRegistryFilterer) ParseRequestCreated(log types.Log) (*BakeRegistryRequestCreated, error) {
	event := new(BakeRegistryRequestCreated)
	if err := _BakeRegistry.contract.UnpackLog(event, "RequestCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BakeRegistrySolutionAcceptedIterator is returned from FilterSolutionAccepted and is used to iterate over the raw logs and unpacked data for SolutionAccepted events raised by the BakeRegistry contract.
type BakeRegistrySolutionAcceptedIterator struct {
	Event *BakeRegistrySolutionAccepted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BakeRegistrySolutionAcceptedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BakeRegistrySolutionAccepted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BakeRegistrySolutionAccepted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BakeRegistrySolutionAcceptedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BakeRegistrySolutionAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BakeRegistrySolutionAccepted represents a SolutionAccepted event raised by the BakeRegistry contract.
type BakeRegistrySolutionAccepted struct {
	RequestId    *big.Int
	Baker        common.Address
	RewardPaid   *big.Int
	RefundAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSolutionAccepted is a free log retrieval operation binding the contract event 0x8cd3d2eafb054005a6c0b9f322ceea8b95889b910e4eb1048f5b5e55dae947b9.
//
// Solidity: event SolutionAccepted(uint256 indexed requestId, address indexed baker, uint256 rewardPaid, uint256 refundAmount)
func (_BakeRegistry *BakeRegistryFilterer) FilterSolutionAccepted(opts *bind.FilterOpts, requestId []*big.Int, baker []common.Address) (*BakeRegistrySolutionAcceptedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var bakerRule []interface{}
	for _, bakerItem := range baker {
		bakerRule = append(bakerRule, bakerItem)
	}

	logs, sub, err := _BakeRegistry.contract.FilterLogs(opts, "SolutionAccepted", requestIdRule, bakerRule)
	if err != nil {
		return nil, err
	}
	return &BakeRegistrySolutionAcceptedIterator{contract: _BakeRegistry.contract, event: "SolutionAccepted", logs: logs, sub: sub}, nil
}

// WatchSolutionAccepted is a free log subscription operation binding the contract event 0x8cd3d2eafb054005a6c0b9f322ceea8b95889b910e4eb1048f5b5e55dae947b9.
//
// Solidity: event SolutionAccepted(uint256 indexed requestId, address indexed baker, uint256 rewardPaid, uint256 refundAmount)
func (_BakeRegistry *BakeRegistryFilterer) WatchSolutionAccepted(opts *bind.WatchOpts, sink chan<- *BakeRegistrySolutionAccepted, requestId []*big.Int, baker []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var bakerRule []interface{}
	for _, bakerItem := range baker {
		bakerRule = append(bakerRule, bakerItem)
	}

	logs, sub, err := _BakeRegistry.contract.WatchLogs(opts, "SolutionAccepted", requestIdRule, bakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BakeRegistrySolutionAccepted)
				if err := _BakeRegistry.contract.UnpackLog(event, "SolutionAccepted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSolutionAccepted is a log parse operation binding the contract event 0x8cd3d2eafb054005a6c0b9f322ceea8b95889b910e4eb1048f5b5e55dae947b9.
//
// Solidity: event SolutionAccepted(uint256 indexed requestId, address indexed baker, uint256 rewardPaid, uint256 refundAmount)
func (_BakeRegistry *BakeRegistryFilterer) ParseSolutionAccepted(log types.Log) (*BakeRegistrySolutionAccepted, error) {
	event := new(BakeRegistrySolutionAccepted)
	if err := _BakeRegistry.contract.UnpackLog(event, "SolutionAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BakeRegistrySolutionRejectedIterator is returned from FilterSolutionRejected and is used to iterate over the raw logs and unpacked data for SolutionRejected events raised by the BakeRegistry contract.
type BakeRegistrySolutionRejectedIterator struct {
	Event *BakeRegistrySolutionRejected // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BakeRegistrySolutionRejectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BakeRegistrySolutionRejected)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BakeRegistrySolutionRejected)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BakeRegistrySolutionRejectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BakeRegistrySolutionRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BakeRegistrySolutionRejected represents a SolutionRejected event raised by the BakeRegistry contract.
type BakeRegistrySolutionRejected struct {
	RequestId *big.Int
	Baker     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSolutionRejected is a free log retrieval operation binding the contract event 0x80137d8518d82ba3a57713b3d1d5aa98e97d2e512d9a2dd2ef75e556163e69b3.
//
// Solidity: event SolutionRejected(uint256 indexed requestId, address indexed baker)
func (_BakeRegistry *BakeRegistryFilterer) FilterSolutionRejected(opts *bind.FilterOpts, requestId []*big.Int, baker []common.Address) (*BakeRegistrySolutionRejectedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var bakerRule []interface{}
	for _, bakerItem := range baker {
		bakerRule = append(bakerRule, bakerItem)
	}

	logs, sub, err := _BakeRegistry.contract.FilterLogs(opts, "SolutionRejected", requestIdRule, bakerRule)
	if err != nil {
		return nil, err
	}
	return &BakeRegistrySolutionRejectedIterator{contract: _BakeRegistry.contract, event: "SolutionRejected", logs: logs, sub: sub}, nil
}

// WatchSolutionRejected is a free log subscription operation binding the contract event 0x80137d8518d82ba3a57713b3d1d5aa98e97d2e512d9a2dd2ef75e556163e69b3.
//
// Solidity: event SolutionRejected(uint256 indexed requestId, address indexed baker)
func (_BakeRegistry *BakeRegistryFilterer) WatchSolutionRejected(opts *bind.WatchOpts, sink chan<- *BakeRegistrySolutionRejected, requestId []*big.Int, baker []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var bakerRule []interface{}
	for _, bakerItem := range baker {
		bakerRule = append(bakerRule, bakerItem)
	}

	logs, sub, err := _BakeRegistry.contract.WatchLogs(opts, "SolutionRejected", requestIdRule, bakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BakeRegistrySolutionRejected)
				if err := _BakeRegistry.contract.UnpackLog(event, "SolutionRejected", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSolutionRejected is a log parse operation binding the contract event 0x80137d8518d82ba3a57713b3d1d5aa98e97d2e512d9a2dd2ef75e556163e69b3.
//
// Solidity: event SolutionRejected(uint256 indexed requestId, address indexed baker)
func (_BakeRegistry *BakeRegistryFilterer) ParseSolutionRejected(log types.Log) (*BakeRegistrySolutionRejected, error) {
	event := new(BakeRegistrySolutionRejected)
	if err := _BakeRegistry.contract.UnpackLog(event, "SolutionRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BakeRegistrySolutionSubmittedIterator is returned from FilterSolutionSubmitted and is used to iterate over the raw logs and unpacked data for SolutionSubmitted events raised by the BakeRegistry contract.
type BakeRegistrySolutionSubmittedIterator struct {
	Event *BakeRegistrySolutionSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BakeRegistrySolutionSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BakeRegistrySolutionSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BakeRegistrySolutionSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BakeRegistrySolutionSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BakeRegistrySolutionSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BakeRegistrySolutionSubmitted represents a SolutionSubmitted event raised by the BakeRegistry contract.
type BakeRegistrySolutionSubmitted struct {
	RequestId    *big.Int
	Baker        common.Address
	SolutionHash [32]byte
	TokenUsage   *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSolutionSubmitted is a free log retrieval operation binding the contract event 0xc30ff78e01d5f02897f9ca8483b38f99495aff4d694eb84cf077e46e38460561.
//
// Solidity: event SolutionSubmitted(uint256 indexed requestId, address indexed baker, bytes32 solutionHash, uint256 tokenUsage)
func (_BakeRegistry *BakeRegistryFilterer) FilterSolutionSubmitted(opts *bind.FilterOpts, requestId []*big.Int, baker []common.Address) (*BakeRegistrySolutionSubmittedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var bakerRule []interface{}
	for _, bakerItem := range baker {
		bakerRule = append(bakerRule, bakerItem)
	}

	logs, sub, err := _BakeRegistry.contract.FilterLogs(opts, "SolutionSubmitted", requestIdRule, bakerRule)
	if err != nil {
		return nil, err
	}
	return &BakeRegistrySolutionSubmittedIterator{contract: _BakeRegistry.contract, event: "SolutionSubmitted", logs: logs, sub: sub}, nil
}

// WatchSolutionSubmitted is a free log subscription operation binding the contract event 0xc30ff78e01d5f02897f9ca8483b38f99495aff4d694eb84cf077e46e38460561.
//
// Solidity: event SolutionSubmitted(uint256 indexed requestId, address indexed baker, bytes32 solutionHash, uint256 tokenUsage)
func (_BakeRegistry *BakeRegistryFilterer) WatchSolutionSubmitted(opts *bind.WatchOpts, sink chan<- *BakeRegistrySolutionSubmitted, requestId []*big.Int, baker []common.Address) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var bakerRule []interface{}
	for _, bakerItem := range baker {
		bakerRule = append(bakerRule, bakerItem)
	}

	logs, sub, err := _BakeRegistry.contract.WatchLogs(opts, "SolutionSubmitted", requestIdRule, bakerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BakeRegistrySolutionSubmitted)
				if err := _BakeRegistry.contract.UnpackLog(event, "SolutionSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSolutionSubmitted is a log parse operation binding the contract event 0xc30ff78e01d5f02897f9ca8483b38f99495aff4d694eb84cf077e46e38460561.
//
// Solidity: event SolutionSubmitted(uint256 indexed requestId, address indexed baker, bytes32 solutionHash, uint256 tokenUsage)
func (_BakeRegistry *BakeRegistryFilterer) ParseSolutionSubmitted(log types.Log) (*BakeRegistrySolutionSubmitted, error) {
	event := new(BakeRegistrySolutionSubmitted)
	if err := _BakeRegistry.contract.UnpackLog(event, "SolutionSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
