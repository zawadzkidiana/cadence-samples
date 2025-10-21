# Dockerized Cadence Samples

## Overview
This project adds Docker support to Cadence samples, making them easier to deploy in containerized environments like Kubernetes.

## Why Dockerize?
- Enable easy deployment to cloud environments
- Ensure consistent runtime environment
- Simplify the getting started experience
- Support Kubernetes deployment scenarios

## Directory Structure
```
samples/docker/
├── Dockerfile              # Main Dockerfile for samples
├── docker-compose.yml      # Local development setup
└── k8s/                   # Kubernetes deployment files
```

## Getting Started
1. Build the Docker image:
   ```bash
   docker build -t cadence-samples .
   ```

2. Run using Docker Compose:
   ```bash
   docker-compose up
   ```

## Testing
- Unit tests: `make test`
- Integration tests: `make integration-test`
- Manual verification steps in [TESTING.md](TESTING.md)

## What's Next?
- [ ] Add Kubernetes Helm charts
- [ ] Add CI/CD pipeline for automated builds
- [ ] Add more sample workflows

# Commit with the clear message
git commit -m "docker: Add hello world workflow sample

Simple containerization of the Cadence hello world workflow example.
Includes:
- Basic worker implementation with workflow and activity
- Simple Dockerfile for building and running the worker
- Environment variable support for Cadence server connection"

# Force push to update your branch
git push -f fork docker-hello-world