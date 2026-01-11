# Project Scanner

A simple CLI tool to scan local directories and list development projects.

## Features

- **Multi-language support**: Detects Go, Node, Python, DotNet, Rust, Java, and Git projects
- **Priority-based detection**: More specific project types take precedence over generic ones
- **Recursive scanning**: Scans directories recursively, stopping when a project is found
- **Configurable paths**: Define which directories to scan via YAML configuration
- **Simple output**: Tab-separated format for easy parsing

## Build Instructions

```bash
# Build the binary
go build -o project-scanner .

# Run the scanner
./project-scanner
```

## Configuration

Edit `config.yaml` to specify the paths you want to scan:

```yaml
paths:
  - /absolute/path/one
  - /absolute/path/two
```

## Project Detection

The tool detects projects by looking for these files/directories in priority order:

### Go (Highest Priority)
- `go.mod`, `go.sum`

### Node
- `package.json`, `package-lock.json`, `yarn.lock`

### Python
- `requirements.txt`, `pyproject.toml`, `setup.py`, `Pipfile`
- Files with `.py` extension

### DotNet
- Files with `.csproj`, `.sln` extensions

### Rust
- `Cargo.toml`, `Cargo.lock`

### Java
- `pom.xml`, `build.gradle`
- Files with `.java` extension

### Git (Lowest Priority)
- `.git` directory

**Note**: Priority is important - if a directory contains both `.git` and `go.mod`, it will be identified as "Go" project, not "Git".

## Output Format

```
<Type>    <Path>
```

Example:
```
Go    /home/user/projects/my-go-app
Node    /home/user/projects/frontend
Python    /home/user/projects/ml-experiment
DotNet    /home/user/projects/web-api
Rust    /home/user/projects/cli-tool
Java    /home/user/projects/spring-app
Git    /home/user/experiments/old-project
```