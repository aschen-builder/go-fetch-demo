# Go `fetch` Demo

This repo includes a small demonstration of a `fetch` CLI tool written in Go.

The current mirroring implementation only supports a shallow mirror.

## Requirements

- Docker
- Git

## Installation

```bash
#!/usr/bin/env bash

git clone https://github.com/aschen-builder/go-fetch-demo.git
cd go-fetch-demo

# Build the Docker image
docker build --rm -t go-fetch:demo .

# Run the Docker image which will open an interactive shell
docker run -it go-fetch:demo

# Run the `fetch` CLI tool
fetch https://www.google.com
fetch --metadata https://www.google.com
fetch --metadata https://www.google.com https://www.amazon.com
fetch --metadata --mirror https://www.google.com
fetch --mirror https://www.google.com

```
