package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	operator "github.com/kubemart/kubemart-operator/api/v1alpha1"
)

func TestMain(m *testing.M) {
	// do stuffs before tests
	fmt.Printf("*********************************************************************\n")
	fmt.Printf("Deleting .env file...\n")
	envFilePath := fmt.Sprintf("%s/.env", scriptsDir)
	_ = os.Remove(envFilePath)

	fmt.Printf("Deleting installation-info.json file...\n")
	infoFilePath := fmt.Sprintf("%s/installation-info.json", scriptsDir)
	_ = os.Remove(infoFilePath)

	fmt.Printf("Deleting all apps (if any) using kubectl and wait for it to finish...\n")
	cmd := exec.Command("kubectl", "delete", "apps", "--force", "--grace-period", "0", "--all", "-A")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
	}

	exitVal := m.Run()
	// do stuffs after tests
	os.Exit(exitVal)
}

func TestGetKubemartApp1(t *testing.T) {
	expectedError := "apps.kubemart.civo.com \"rabbitmq\" not found"
	_, actualError := GetKubemartApp("rabbitmq", "kubemart-system")
	if expectedError != actualError.Error() {
		t.Errorf("Expected %s but actual is %s", expectedError, actualError)
	}
}

func TestGetKubemartApp2(t *testing.T) {
	manifest := `
apiVersion: kubemart.civo.com/v1alpha1
kind: App
metadata:
  name: rabbitmq
  namespace: kubemart-system
spec:
  name: rabbitmq
  action: install`

	fileBytes := []byte(manifest)
	filename := "./rabbitmq.yaml"
	_ = ioutil.WriteFile(filename, fileBytes, 0644)

	c1 := exec.Command("cat", filename)
	c2 := exec.Command("kubectl", "apply", "-f", "-")
	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = os.Stdout
	_ = c2.Start()
	_ = c1.Run()
	_ = c2.Wait()
	time.Sleep(1 * time.Second)

	_, actualError := GetKubemartApp("rabbitmq", "kubemart-system")
	if actualError != nil {
		t.Errorf("Expected nil error but actual is %s", actualError)
	}
}

func TestBase64Decode(t *testing.T) {
	helloEncoded := "aGVsbG8="
	expected := "hello"
	actual, _ := Base64Decode(helloEncoded)
	if expected != actual {
		t.Errorf("Expected %s but actual is %s", expected, actual)
	}
}

func TestCreateEnvFileFromConfig(t *testing.T) {
	configs := []operator.Configuration{
		{
			Key:           "KEY_1",
			Value:         "VALUE_1",
			ValueIsBase64: false,
		},
		{
			Key:           "KEY_2",
			Value:         "VkFMVUVfMg==", // "VALUE_2" in base64
			ValueIsBase64: true,
		},
	}

	_ = CreateEnvFileFromConfig(configs)
	envFilePath := fmt.Sprintf("%s/.env", scriptsDir)
	fileBytes, _ := ioutil.ReadFile(envFilePath)
	actual := string(fileBytes)

	expectedLine1 := `export KEY_1="VALUE_1"`
	if !strings.Contains(actual, expectedLine1) {
		t.Errorf("Expected %s from line 1 but actual file content is %s", expectedLine1, actual)
	}

	expectedLine2 := `export KEY_2="VALUE_2"`
	if !strings.Contains(actual, expectedLine2) {
		t.Errorf("Expected %s from line 2 but actual file content is %s", expectedLine2, actual)
	}
}

func TestSaveInstallationInfo(t *testing.T) {
	appName := "wordpress"
	appNamespace := "kubemart-system"
	_ = SaveInstallationInfo(appName, appNamespace)
	infoFilePath := fmt.Sprintf("%s/installation-info.json", scriptsDir)

	file, _ := ioutil.ReadFile(infoFilePath)
	iif := InstallationInfoFile{}
	_ = json.Unmarshal([]byte(file), &iif)

	expectedAppName := appName
	actualAppName := iif.Name
	if expectedAppName != actualAppName {
		t.Errorf("Expected %s but actual is %s", expectedAppName, actualAppName)
	}

	expectedAppNamespace := appName
	actualAppNamespace := iif.Name
	if expectedAppNamespace != actualAppNamespace {
		t.Errorf("Expected %s but actual is %s", expectedAppNamespace, actualAppNamespace)
	}
}
