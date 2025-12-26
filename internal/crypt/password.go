package crypt

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

// PromptPassword asks user for password without echoing
func PromptPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	// Read password without echo
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	fmt.Println() // newline after password entry
	return string(bytePassword), nil
}

// PromptPasswordConfirm asks twice and verifies they match
func PromptPasswordConfirm(prompt string) (string, error) {
	password, err := PromptPassword(prompt)
	if err != nil {
		return "", err
	}

	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	confirm, err := PromptPassword("Confirm password: ")
	if err != nil {
		return "", err
	}

	if password != confirm {
		return "", fmt.Errorf("passwords do not match")
	}

	return password, nil
}
