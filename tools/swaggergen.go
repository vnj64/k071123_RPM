package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// serviceDirs holds the list of service directories relative to the repository root.
var serviceDirs = []string{
	"internal/services/notification_service",
	"internal/services/parking_service",
	"internal/services/user_service",
	"internal/services/order_service",
}

// swagFileRelativePath is the relative location of the generated swagger file within each service directory.
const swagFileRelativePath = "docs/swagger.json"

// docsFolderName is the folder that will be deleted after generation.
const docsFolderName = "docs"

// outputDir is where each service’s swagger JSON file will be copied.
const outputDir = "docs_aggregator/public"

func main() {
	// Ensure the output directory exists.
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory %s: %v", outputDir, err)
	}

	// Process each service directory.
	for _, dir := range serviceDirs {
		serviceName := filepath.Base(dir)
		fmt.Printf("Processing service: %s\n", serviceName)

		// 1. Run "swag init" in the service directory with the --overridesFile flag.
		cmd := exec.Command("swag", "init", "-g", filepath.Join("cmd", "main.go"), "--parseDependency", "--parseInternal")
		cmd.Dir = dir

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error running swag init in %s: %v\nStderr: %s", dir, err, stderr.String())
			continue
		}
		fmt.Printf("Swag init output for %s:\n%s\n", dir, stdout.String())

		// 2. Read the generated swagger.json file.
		swagPath := filepath.Join(dir, swagFileRelativePath)
		data, err := ioutil.ReadFile(swagPath)
		if err != nil {
			log.Printf("Error reading swagger file for %s at %s: %v", serviceName, swagPath, err)
			continue
		}

		// 3. Save the file into docs_aggregator/public with the service name.
		destFile := filepath.Join(outputDir, serviceName+".json")
		if err := ioutil.WriteFile(destFile, data, 0644); err != nil {
			log.Printf("Error writing swagger file for %s to %s: %v", serviceName, destFile, err)
			continue
		}
		fmt.Printf("Saved %s swagger JSON to %s\n", serviceName, destFile)

		// 4. Delete the local docs folder in the service directory.
		docsPath := filepath.Join(dir, docsFolderName)
		if err := os.RemoveAll(docsPath); err != nil {
			log.Printf("Warning: could not delete docs folder for %s at %s: %v", serviceName, docsPath, err)
		} else {
			fmt.Printf("Deleted local docs folder for %s\n", serviceName)
		}
	}

	// Optionally, list the files written.
	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		log.Fatalf("Error listing files in output directory: %v", err)
	}
	var fileNames []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			fileNames = append(fileNames, f.Name())
		}
	}
	fmt.Printf("Swagger files saved: %v\n", fileNames)

	// 5. Start a simple HTTP server to serve the docs_aggregator/public folder.
	fs := http.FileServer(http.Dir(outputDir))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	port := ":8129"
	fmt.Printf("Serving Swagger files at http://localhost%s/swagger/\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
