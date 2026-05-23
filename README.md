<!-- # Terraform_DevOps



## Getting Started
**Terraform_runner** Project for automating and Integrating the Terraform projects with CI/CD pipelines and also for making the process of using Terraform easier and faster for developers and DevOps engineers , also with vCenter (VMware) and for managing infrastructure as code in a more efficient way.

Usage : 
* **first** clone the repository of project to your local space : 

    ```sh
    git clone https://github.com/Ali-Raven/Terraform_Runner.git
    ```
* **build** the executable file to use simply : 
    ```sh
    go build .
    ```
* going to the direcoty of your project and then **change the code with your idea :**

    ```sh
    cd Terraform_Runner/
     
    # use the executable
    ./terraform_runner <command>
    ``` -->

# Terraform DevOps CLI

A Go-based CLI project for DevOps workflow automation, including Terraform, Helm, YAML, staging, web UI, and server management helpers.

> [!NOTE]
> This README was generated from the repository layout and file names. For exact command usage, verify the CLI flags in `main.go` and related command files.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Overview

This repository provides a CLI tool for managing Infrastructure as Code and DevOps workflows. It appears to include support for Terraform-related automation, Helm chart management, YAML processing, staging environments, and an optional web UI layer.

## Features

- CLI entrypoint in `main.go`
- Helm integration in `helm.go`
- YAML configuration support in `yml.go`
- Staging management in `stage.go`
- Web UI support in `webui.go`
- Custom workflow modules in `cyborg.go`, `nozaros.go`, and `oranos.go`
- Shared utilities under `helper/`

> [!WARNING]
> Make sure you have the proper environment configured before running commands. Missing tools or credentials may prevent successful execution.

## Project Structure

- `main.go` — CLI bootstrap and command routing
- `cyborg.go` — custom DevOps workflow or command logic
- `helm.go` — Helm-related functionality
- `nozaros.go` — workflow module or environment operations
- `oranos.go` — workflow module or environment operations
- `stage.go` — stage/environment management
- `webui.go` — web interface launcher or server
- `yml.go` — YAML generation and parsing
- `helper/` — helper utilities and shared functions
- `servers/` — server-related configuration or deployment assets
- `terraform_DevOps` — compiled binary
- `LICENSE` — project license
- `.gitignore` — ignored files

## Requirements

- Go installed (version compatible with `go.mod`)
- Terraform installed if using Terraform-related commands
- Helm installed if using Helm-related commands
- Any required cloud or cluster credentials configured

> [!TIP]
> Run `go env` and `go version` to confirm your Go environment.

## Installation

1. Clone the repository:

   ```bash
   git clone /home/raven/Desktop/Workstation/cli_devops_project/terraform
   cd /home/raven/Desktop/Workstation/cli_devops_project/terraform
   ```

2. Build the CLI:

   ```bash
   go build -o terraform_DevOps .
   ```

3. Run the binary:

   ```bash
   ./terraform_DevOps --help
   ```

## Usage

Use the built CLI binary to inspect available commands and run tasks.

```bash
./terraform_DevOps --help
```

Common commands may include:

```bash
./terraform_DevOps deploy
./terraform_DevOps helm install
./terraform_DevOps stage create
./terraform_DevOps webui start
```

> [!NOTE]
> Replace the example commands above with the actual command names defined in `main.go` and related files.

## Development

To modify the project:

1. Open the code in your editor.
2. Update files in the repository.
3. Rebuild the binary after changes:

   ```bash
   go build -o terraform_DevOps .
   ```

4. Run the CLI to test your changes.

## Testing

If unit tests are added later, use:

```bash
go test ./...
```

> [!IMPORTANT]
> There are no test files listed in the repository currently. Add tests under the relevant packages as needed.

## Contributing

- Fork the repository
- Create a new branch for your feature or fix
- Make changes and test locally
- Submit a pull request

## License

This project is licensed under the terms in the `LICENSE` file.

> [!NOTE]
> Verify the license contents in `LICENSE` before redistributing or modifying.


