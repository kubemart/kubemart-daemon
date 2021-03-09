package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	operator "github.com/kubemart/kubemart-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	scriptsDir = "../../scripts"
)

// InstallationInfoFile ...
type InstallationInfoFile struct {
	Name      string `json:"cr_name"`
	Namespace string `json:"cr_namespace"`
}

// GetKubemartApp ...
func GetKubemartApp(name, namespace string) (*operator.App, error) {
	app := &operator.App{}

	config, err := rest.InClusterConfig()
	if err != nil {
		return app, err
	}

	scheme := runtime.NewScheme()
	operator.AddToScheme(scheme)
	operatorClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return app, err
	}

	target := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	err = operatorClient.Get(context.Background(), target, app)
	if err != nil {
		return app, err
	}

	return app, nil
}

// Base64Decode takes base64-encoded string and returns its original value
func Base64Decode(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	decoded := string(data)
	return decoded, nil
}

// AppendEnvFile ...
func AppendEnvFile(textToAppend string) error {
	filepath := fmt.Sprintf("%s/.env", scriptsDir)
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	text := fmt.Sprintf("%s\n", textToAppend)
	_, err = f.WriteString(text)
	if err != nil {
		return err
	}

	return nil
}

// CreateEnvFileFromConfig ...
func CreateEnvFileFromConfig(configs []operator.Configuration) error {
	for _, config := range configs {
		var value string

		if config.ValueIsBase64 {
			decoded, err := Base64Decode(config.Value)
			if err != nil {
				return err
			}
			value = decoded
		} else {
			value = config.Value
		}

		envStr := fmt.Sprintf("export %s=\"%s\"", config.Key, value)
		err := AppendEnvFile(envStr)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveInstallationInfo ...
func SaveInstallationInfo(appName string, namespace string) error {
	jsonFile := InstallationInfoFile{}
	jsonFile.Name = appName
	jsonFile.Namespace = namespace

	file, _ := json.MarshalIndent(jsonFile, "", " ")
	filepath := fmt.Sprintf("%s/installation-info.json", scriptsDir)
	err := ioutil.WriteFile(filepath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
