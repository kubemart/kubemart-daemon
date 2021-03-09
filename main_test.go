package main

import (
	"testing"
)

func TestParseFlags(t *testing.T) {
	expectedError := "--app-name flag is empty"
	_, actualError := parseFlags()
	if expectedError != actualError.Error() {
		t.Errorf("Expected %s but actual is %s", expectedError, actualError)
	}
}
