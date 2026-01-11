package main

import (
	"fmt"
	"log"

	"project-scanner/scanner"
)

func main() {
	// Read configuration
	config, err := scanner.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Scan for projects
	projects, err := scanner.ScanProjects(config.Paths)
	if err != nil {
		log.Fatalf("Failed to scan projects: %v", err)
	}

	// Output results
	for _, project := range projects {
		fmt.Printf("%s\t%s\n", project.Type, project.Path)
	}
}
