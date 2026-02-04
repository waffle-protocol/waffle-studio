package ai

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetAPIKeyPath returns the path to the Gemini API key file
func GetAPIKeyPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".waffle", "gemini_apikey"), nil
}

// GetAPIKey reads the Gemini API key from file or prompts the user to enter it
func GetAPIKey() (string, error) {
	keyPath, err := GetAPIKeyPath()
	if err != nil {
		return "", err
	}

	// Check environment variable
	if envKey := os.Getenv("GEMINI_API_KEY"); envKey != "" {
		return envKey, nil
	}

	// Try to read existing key
	data, err := os.ReadFile(keyPath)
	if err == nil {
		key := strings.TrimSpace(string(data))
		if key != "" {
			return key, nil
		}
	}

	// Prompt user for API key
	fmt.Println()
	fmt.Println("Gemini API key not found.")
	fmt.Println("Get your API key from: https://aistudio.google.com/app/apikey")
	fmt.Print("Enter your Gemini API key: ")

	reader := bufio.NewReader(os.Stdin)
	key, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}
	key = strings.TrimSpace(key)

	if key == "" {
		return "", fmt.Errorf("API key cannot be empty")
	}

	// Save the key for future use
	if err := SaveAPIKey(key); err != nil {
		fmt.Printf("Warning: failed to save API key: %v\n", err)
	}

	return key, nil
}

// SaveAPIKey saves the Gemini API key to file
func SaveAPIKey(key string) error {
	keyPath, err := GetAPIKeyPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(keyPath), 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write key with restricted permissions
	if err := os.WriteFile(keyPath, []byte(key), 0600); err != nil {
		return fmt.Errorf("failed to write API key: %w", err)
	}

	return nil
}
