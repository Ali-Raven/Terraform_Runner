# Terraform DevOps CLI
<p align="center">
  <img src="https://github.com/user-attachments/assets/b8b896a6-19d0-46ae-8ca2-62ab81e016a7" alt="Go App Demo" style="width: 100%; max-width: 900px;" />
</p>

A Go-based CLI project for DevOps workflow automation, including Terraform, Helm, YAML, staging, web UI, and server management helpers.

> [!NOTE]
> This README was generated from the repository layout and file names. For exact command usage, verify the CLI flags in `main.go` and related command files.

## Table of Contents

- [**Overview**](#overview)
- [**Project Structure**](#project-structure)
- [**Requirements**](#requirements)
- [**Installation**](#installation)
- [**Usage**](#usage)
- [**Development**](#development)
- [**Contributing**](#contributing)
- [**License**](#license)

## Overview

This repository provides a CLI tool for managing Infrastructure as Code and DevOps workflows. It appears to include support for **Terraform-related automation**, **Helm (Esxi and vCenter)** management, **working with Ansible** and an optional **web UI** layer.

> [!WARNING]
> Make sure you have the proper environment configured before running commands. Missing tools or credentials may prevent successful execution.

<!-- ## Project Structure

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
- `.gitignore` — ignored files -->

## Requirements

- Go installed (version compatible with `go.mod`)
- **Terraform** installed if using Terraform-related commands
- Any required cloud or cluster credentials configured

> [!TIP]
> Run `go env` and `go version` to confirm your Go environment.

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/Ali-Raven/Terraform_Runner.git
   cd Terraform_Runner
   ```

2. **Build the CLI**:

   ```bash
   go build .
   ```
   or if you want specific name for you **executable binary**
   ```bash
   go build -o <YOUR BINARY NAME> <PATH>
   ```

3. **Run the binary**:

   ```bash
   ./<Your binary name> [flags]
   ```

## Usage

Use the built CLI binary to inspect available commands and run tasks.

```bash
./terraform_DevOps --help
```

Common commands may include:

```bash
./terraform_DevOps --nozaros
./terraform_DevOps --helm
./terraform_DevOps --oranos
./terraform_DevOps --cyborg
```

## Development

To modify the project:

1. Open the code in your editor.
2. Update files in the repository.
3. Rebuild the binary after changes:

   ```bash
   go build -o terraform_DevOps .
   ```

4. Run the CLI to test your changes.

## Contributing

- **Fork** the repository
- Create a **new branch for your feature or fix**
- Make changes and **test locally**
- **Submit a pull request**

## License

**MIT**


