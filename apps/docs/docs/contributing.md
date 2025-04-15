---
id: contributing
title: Contributing & Development
---

# Contributing & Development

We welcome contributions to the ShareKube project! This guide will help you set up a development environment and understand how to contribute to the project.

## Development Environment Setup

### Prerequisites

- Go 1.16 or higher
- Kubernetes cluster for testing (minikube, kind, or k3d work well)
- Docker for building images
- kubectl configured for your cluster
- Operator SDK (optional, for operator development)

### Getting the Source Code

Clone the repository:

```bash
git clone https://github.com/miloszsobczak/sharekube.git
cd sharekube
```

### Building from Source

Build the operator:

```bash
# Build the operator binary
make build

# Build the Docker image
make docker-build IMG=sharekube/operator:dev
```

### Running Locally

Run the operator locally:

```bash
# Install CRDs
make install

# Run the operator
make run
```

### Testing

Run tests:

```bash
# Run unit tests
make test

# Run integration tests
make test-integration
```

## Project Structure

The KubeShare project follows a standard Kubernetes operator structure:

```
├── api/                  # API definitions (CRDs)
│   └── v1alpha1/         # API version
├── config/               # Operator configuration
├── controllers/          # Operator controllers
├── pkg/                  # Shared packages
│   ├── resources/        # Resource management
│   └── transformation/   # Transformation logic (future)
└── test/                 # Test files
```

## Development Workflow

1. **Fork the Repository**: Fork the ShareKube repository on GitHub
2. **Create a Branch**: Create a branch for your feature or bug fix
3. **Make Changes**: Develop and test your changes
4. **Write Tests**: Add tests for your changes
5. **Submit a Pull Request**: Open a PR with a description of your changes

## Coding Standards

### Go Code

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Format code with `gofmt`
- Document all exported functions, types, and constants
- Write unit tests for all functionality

### Kubernetes Resources

- Follow the [Kubernetes API Conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
- Use Kubernetes style for YAML indentation (2 spaces)

## Adding a New Feature

1. **Design**: Document the feature design and discuss in an issue
2. **API Changes**: If changing the API, update the CRD definitions
3. **Implementation**: Implement the feature with tests
4. **Documentation**: Update documentation to reflect changes

## Future Transformation Rules Development

When implementing transformation rules (planned future feature):

1. Define the transformation in the `pkg/transformation` package
2. Create a transformation handler for each resource type
3. Add unit tests for transformations
4. Update the API documentation

## Release Process

1. **Version Bump**: Update version numbers in the code and manifests
2. **Changelog**: Update the CHANGELOG.md with a summary of changes
3. **Tag Release**: Create a Git tag for the release
4. **CI/CD**: The CI pipeline will build and publish artifacts

## Community

Join our community:

- **GitHub Discussions**: For feature discussions and community help
- **Issues**: Bug reports and feature requests
- **Pull Requests**: Code contributions

## Code of Conduct

We follow the CNCF Code of Conduct. Please be respectful and collaborative in all interactions.

## License

KubeShare is licensed under the Apache 2.0 License. All contributions must comply with this license.

---

Thank you for considering a contribution to KubeShare! Your help makes the project better for everyone. 