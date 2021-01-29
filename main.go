/*
Useful resources
- https://blog.rapid7.com/2016/08/04/build-a-simple-cli-tool-with-golang/
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// HelmChartHeader is used to extract metadata.name from HelmChart YAML files
type HelmChartHeader struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}

// InstallationInfo is used to save JSON file into "scripts" folder
type InstallationInfo struct {
	AppName          string
	CrdName          string
	Namespace        string
	HelmReleaseName  string
	HelmMetadataName string
}

var isHelmChart bool = false
var outputDir string = "./scripts"
var helmChartFileName string = "helm-chart.yaml"
var helmChartFilePath string = fmt.Sprintf("%s/%s", outputDir, helmChartFileName)

func main() {
	appName := flag.String("app-name", "", "Marketplace App name (required)")
	crdAppName := flag.String("crd-app-name", "", "The metadata.name of the App CRD (required)")
	namespaceName := flag.String("namespace", "default", "Namespace is the namespace where the App CR lives (will use default if it's empty)")
	flag.Parse()

	// Check CLI flags
	if *appName == "" || *crdAppName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Printf("App name: %s, CRD App name: %s, Namespace: %s\n", *appName, *crdAppName, *namespaceName)

	// Return if marketplace app doesn't have app.yaml file
	yamlFilepath := fmt.Sprintf("./marketplace/%s/app.yaml", *appName)
	if !fileExists(yamlFilepath) {
		log.Println("This app doesn't have app.yaml file, will proceed with install.sh file using bash")
		// Print installation information to JSON file
		err := saveInstallationInfo(isHelmChart, *appName, *crdAppName, *namespaceName)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Open file
	appFile := fmt.Sprintf("./marketplace/%s/app.yaml", *appName)
	b, err := ioutil.ReadFile(appFile)

	if err != nil {
		log.Fatal(err)
	}

	// Get file content
	content := string(b) // may contain line breaks
	// fmt.Println("Original")
	// fmt.Println(content)

	// Split YAMLs
	counter := 1
	splitted := strings.Split(content, "---")
	for _, s := range splitted {
		cleaned := strings.TrimSuffix(strings.TrimPrefix(s, "\n"), "\n")
		if cleaned == "" {
			continue // skip the empty one
		}
		// fmt.Println("Splitted")
		// fmt.Printf("%s\n\n", cleaned)

		// Create parent directory if it's not exists
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			err = os.Mkdir(outputDir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Save file to disk
		if strings.Contains(cleaned, "HelmChart") {
			err = saveFile(cleaned, helmChartFilePath)
			if err != nil {
				log.Fatalf("Failed to save file %s", helmChartFilePath)
			}
			isHelmChart = true
		} else {
			filename := fmt.Sprintf("%s/k8s-%d.yaml", outputDir, counter)
			err = saveFile(cleaned, filename)
			if err != nil {
				log.Fatalf("Failed to save file %s", filename)
			}
			counter++
		}
	}

	// Print installation information to JSON file
	err = saveInstallationInfo(isHelmChart, *appName, *crdAppName, *namespaceName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Go app part is done")
}

// Return true if file exists at path
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Save HelmChart info to JSON file
func saveInstallationInfo(isHelmChart bool, appName string, crdAppName string, namespace string) error {
	jsonFile := InstallationInfo{}
	jsonFile.AppName = appName
	jsonFile.CrdName = crdAppName
	jsonFile.Namespace = namespace

	if isHelmChart {
		hch, err := readChartHeader()
		if err != nil {
			return err
		}

		jsonFile.HelmReleaseName = hch.Metadata.Name
		jsonFile.HelmMetadataName = hch.Metadata.Name
	}

	file, _ := json.MarshalIndent(jsonFile, "", " ")
	filepath := fmt.Sprintf("%s/installation-info.json", outputDir)
	err := ioutil.WriteFile(filepath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Read HelmChart YAML file and fill the HelmChartHeader struct
func readChartHeader() (*HelmChartHeader, error) {
	buf, err := ioutil.ReadFile(helmChartFilePath)
	if err != nil {
		return nil, err
	}

	hch := &HelmChartHeader{}
	err = yaml.Unmarshal(buf, hch)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", helmChartFilePath, err)
	}

	return hch, nil
}

// Save individual YAML file
func saveFile(content string, filename string) error {
	data := []byte(content)
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
