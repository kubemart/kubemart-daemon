/*
Useful resources
- https://blog.rapid7.com/2016/08/04/build-a-simple-cli-tool-with-golang/
*/

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kubemart/kubemart-daemon/pkg/utils"
)

type cliFlags struct {
	appName      string
	appNamespace string
}

func main() {
	cf, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("App name: %s, Namespace: %s\n", cf.appName, cf.appNamespace)

	app, err := utils.GetKubemartApp(cf.appName, cf.appNamespace)
	if err != nil {
		log.Fatal(err)
	}

	configs := app.Status.Configurations
	err = utils.CreateEnvFileFromConfig(configs)
	if err != nil {
		log.Fatal(err)
	}

	err = utils.SaveInstallationInfo(cf.appName, cf.appNamespace)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Go app part is done")
}

func parseFlags() (cliFlags, error) {
	cf := cliFlags{}
	appName := flag.String("app-name", "", "Marketplace App name (required)")
	namespaceName := flag.String("namespace", "kubemart-system", "Namespace is the namespace where the App CR lives (will use kubemart-system if it's empty)")
	flag.Parse()

	if *appName == "" {
		return cf, fmt.Errorf("--app-name flag is empty")
	}

	cf.appName = *appName
	cf.appNamespace = *namespaceName
	return cf, nil
}
