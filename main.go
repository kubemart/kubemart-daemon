/*
Useful resources
- https://blog.rapid7.com/2016/08/04/build-a-simple-cli-tool-with-golang/
*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/civo/bizaar-daemon/pkg/utils"
	operator "github.com/civo/bizaar-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type installationInfoFile struct {
	Name      string `json:"cr_name"`
	Namespace string `json:"cr_namespace"`
}

func main() {
	appName := flag.String("app-name", "", "Marketplace App name (required)")
	namespaceName := flag.String("namespace", "bizaar-system", "Namespace is the namespace where the App CR lives (will use bizaar-system if it's empty)")
	flag.Parse()

	// Check CLI flags
	if *appName == "" {
		flag.PrintDefaults()
		log.Fatalln("Error: --app-name flag is empty")
	}
	log.Printf("App name: %s, Namespace: %s\n", *appName, *namespaceName)

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	scheme := runtime.NewScheme()
	operator.AddToScheme(scheme)
	operatorClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatal(err)
	}

	target := types.NamespacedName{
		Namespace: *namespaceName,
		Name:      *appName,
	}
	app := &operator.App{}
	err = operatorClient.Get(context.Background(), target, app)
	if err != nil {
		log.Fatal(err)
	}

	configs := app.Status.Configurations
	for _, config := range configs {
		var value string

		if config.ValueIsBase64 {
			value, err = utils.Base64Decode(config.Value)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			value = config.Value
		}

		envStr := fmt.Sprintf("export %s=\"%s\"", config.Key, value)
		appendEnvFile(envStr)
	}

	err = saveInstallationInfo(*appName, *namespaceName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Go app part is done")
}

func saveInstallationInfo(appName string, namespace string) error {
	jsonFile := installationInfoFile{}
	jsonFile.Name = appName
	jsonFile.Namespace = namespace

	file, _ := json.MarshalIndent(jsonFile, "", " ")
	filepath := "./scripts/installation-info.json"
	err := ioutil.WriteFile(filepath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func appendEnvFile(textToAppend string) {
	filepath := "./scripts/.env"
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	text := fmt.Sprintf("%s\n", textToAppend)
	_, err = f.WriteString(text)
	if err != nil {
		log.Fatal(err)
	}
}
