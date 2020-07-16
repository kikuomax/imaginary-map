package common

import (
	"errors"
	"fmt"
	"os"
)

// Returns the value of a given environment variable.
// Returns an empty string and an error if no value is set.
func GetEnv (name string) (string, error) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return "", errors.New(
			fmt.Sprintf(
				"environment variable %s is not set",
				name))
	}
	return value, nil
}
