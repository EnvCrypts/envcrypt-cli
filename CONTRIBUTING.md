# Contributing to EnvCrypt CLI

We love your input! We want to make contributing to EnvCrypt CLI as easy and transparent as possible.

## Development Setup

1. Make sure you have Go 1.22+ installed.
2. Clone the repository:
   ```bash
   git clone https://github.com/envcrypts/envcrypt-cli.git
   cd envcrypt-cli
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build the project:
   ```bash
   go build -o envcrypt main.go
   ```

## Pull Requests

1. Fork the repo and create your branch from `main`.
2. Make sure your code formats well (`go fmt ./...`).
3. Issue that pull request!

## Reporting Bugs

Please include:
* Your operating system and version.
* The version of EnvCrypt CLI (`envcrypt version`).
* Steps to reproduce the bug.
