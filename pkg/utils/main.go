package utils

import "os"

// ReadEnvOrDefault returns the value of the environment variable specified by 'key',
// or the 'defaultValue' if the environment variable is not set.
//
// Parameters:
// - key: the name of the environment variable to retrieve.
// - defaultValue: the default value to return if the environment variable is not set.
//
// Return type:
// - string: the value of the environment variable or the default value.
func ReadEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetOrDefault returns the value if it is not empty, otherwise it returns the defaultValue.
//
// Parameters:
// - value: the value to check if empty.
// - defaultValue: the default value to return if value is empty.
//
// Return type: string
func GetOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
