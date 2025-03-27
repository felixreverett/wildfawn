package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Place for code to import settings data
func LoadSecrets(filename string) (*Secrets, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets file: %v\nDoes the file exist on your local machine?", err)
	}

	var secrets Secrets
	if err := json.Unmarshal(data, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &secrets, nil
}
