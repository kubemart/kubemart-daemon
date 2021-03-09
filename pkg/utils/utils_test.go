package utils

import (
	"testing"
)

func TestGetKubemartApp(t *testing.T) {
	expectedError := "apps.kubemart.civo.com \"rabbitmq\" not found"
	_, actualError := GetKubemartApp("rabbitmq", "kubemart-system")
	if expectedError != actualError.Error() {
		t.Errorf("Expected %s but actual is %s", expectedError, actualError)
	}
}
