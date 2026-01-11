package scanner

import (
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ProjectType defines the structure for project detection rules
type ProjectType struct {
	Name        string
	Files       []string // exact filenames
	Extensions  []string // file extensions (with dot)
	Directories []string // directory names
}

// ProjectTypes defines all supported project types in priority order
var ProjectTypes = []ProjectType{
	{
		Name:  "Go",
		Files: []string{"go.mod", "go.sum"},
	},
	{
		Name:  "Node",
		Files: []string{"package.json", "package-lock.json", "yarn.lock"},
	},
	{
		Name:       "Python",
		Files:      []string{"requirements.txt", "pyproject.toml", "setup.py", "Pipfile"},
		Extensions: []string{".py"},
	},
	{
		Name:       "DotNet",
		Extensions: []string{".csproj", ".sln"},
	},
	{
		Name:  "Rust",
		Files: []string{"Cargo.toml", "Cargo.lock"},
	},
	{
		Name:       "Java",
		Files:      []string{"pom.xml", "build.gradle"},
		Extensions: []string{".java"},
	},
	{
		Name:        "Git",
		Directories: []string{".git"},
	},
}

// Config represents the YAML configuration structure
type Config struct {
	Paths []string `yaml:"paths"`
}

// Project represents a detected development project
type Project struct {
	Type string
	Path string
}

// LoadConfig reads and parses the YAML configuration file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// detectProjectType determines the project type based on files in the directory
// Checks project types in priority order defined in ProjectTypes
func detectProjectType(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	// Build sets of found files, extensions, and directories
	foundFiles := make(map[string]bool)
	foundExtensions := make(map[string]bool)
	foundDirectories := make(map[string]bool)

	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() {
			foundFiles[name] = true
			ext := filepath.Ext(name)
			if ext != "" {
				foundExtensions[ext] = true
			}
		} else {
			foundDirectories[name] = true
		}
	}

	// Check each project type in priority order
	for _, projType := range ProjectTypes {
		if matchesProjectType(projType, foundFiles, foundExtensions, foundDirectories) {
			return projType.Name
		}
	}

	return ""
}

// matchesProjectType checks if a directory matches a specific project type
func matchesProjectType(projType ProjectType, files, extensions, directories map[string]bool) bool {
	// Check exact files
	for _, file := range projType.Files {
		if files[file] {
			return true
		}
	}

	// Check extensions
	for _, ext := range projType.Extensions {
		if extensions[ext] {
			return true
		}
	}

	// Check directories
	for _, dir := range projType.Directories {
		if directories[dir] {
			return true
		}
	}

	return false
}

// ScanProjects recursively scans the given paths for development projects
func ScanProjects(paths []string) ([]Project, error) {
	var projects []Project

	for _, rootPath := range paths {
		err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// Log error but continue scanning
				return nil
			}

			if !d.IsDir() {
				return nil
			}

			// Check if this directory is a project
			projectType := detectProjectType(path)
			if projectType != "" {
				projects = append(projects, Project{
					Type: projectType,
					Path: path,
				})
				// Stop descending into this directory since it's a project
				return filepath.SkipDir
			}

			return nil
		})

		if err != nil {
			// Log error but continue with other paths
			continue
		}
	}

	return projects, nil
}
