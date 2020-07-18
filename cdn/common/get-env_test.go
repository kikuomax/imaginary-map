package common

import (
	"os"
	"testing"
)

const existingVariableName = "THIS_IS_A_TEST_VARIABLE"
const existingVariableValue = "this is a test value"

const nonExistingVariableName = "THIS_VARIABLE_DOES_NOT_EXIST"

// Sets up the test environment.
func setup () error {
	return os.Setenv(existingVariableName, existingVariableValue)
}

// Tears down the test environment.
func teardown () error {
	return os.Unsetenv(existingVariableName)
}

func TestGetEnv (t *testing.T) {
	// setup and deferred teardown
	if err := setup(); err != nil {
		t.Fatal(err)
	}
	defer teardown()
	// runs tests
	t.Run(
		"GetEnv should return a value for an existing environment variable",
		func (t *testing.T) {
			v, err := GetEnv(existingVariableName)
			if v != existingVariableValue {
				t.Errorf(
					"expected GetEnv('%v') == '%v' but got '%v'",
					existingVariableName,
					existingVariableValue,
					v)
			}
			if err != nil {
				t.Errorf(
					"expected GetEnv('%v') returns no error but got %v",
					existingVariableName,
					err)
			}
		})
	t.Run(
		"GetEnv should return an error for a non-existing environment variable",
		func (t *testing.T) {
			v, err := GetEnv(nonExistingVariableName)
			if err == nil {
				t.Errorf(
					"expected GetEnv('%v') returns an error but got nil",
					nonExistingVariableName)
			}
			if v != "" {
				t.Errorf(
					"expected GetEnv('%v') returns an empty string but got '%v'",
					nonExistingVariableName,
					v)
			}
		})
}
