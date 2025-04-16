# Contributing to ShareKube

We welcome contributions to the ShareKube project! This guide will help you set up a development environment and understand how to contribute to the project.

## Related Documents

- [README.md](README.md): Project overview and main documentation
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md): Our community guidelines
- [SECURITY.md](SECURITY.md): Security policy and vulnerability reporting
- [LICENSE](LICENSE): Apache 2.0 license details

## Development Environment Setup

### Prerequisites

- Go 1.19 or higher
- Kubernetes cluster for testing (minikube, kind, or k3d work well)
- Docker for building images
- kubectl configured for your cluster
- Operator SDK (optional, for operator development)
- Node.js 16+ (for documentation site)
- Yarn (for package management)

### Getting the Source Code

Clone the repository:

```bash
git clone https://github.com/miloszsobczak/sharekube.git
cd sharekube
```

### Local Development Workflow

This section provides detailed instructions for setting up a local development environment to build, test, and run ShareKube.

#### Building the Operator

```bash
cd packages/operator
export PATH=$PATH:/usr/local/go/bin  # Ensure Go is in your PATH
go build -o bin/manager main.go
```

#### Setting Up a Local Kubernetes Cluster

For testing the ShareKube operator, you can set up a local Kubernetes cluster using kind:

```bash
# Install kind if not already installed
brew install kind     # macOS with Homebrew
# or
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-$(uname)-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# Create a Kubernetes cluster
kind create cluster --name sharekube

# Save the kubeconfig to a file for easier access
kind get kubeconfig --name=sharekube > /tmp/kubeconfig
```

#### Deploying the CRD

```bash
# Apply the CRD definition
kubectl --kubeconfig=/tmp/kubeconfig apply -f config/crd/bases/sharekube.dev_sharekubes.yaml
```

#### Running the Operator

You can run the operator directly from your development environment:

```bash
./bin/manager --kubeconfig=/tmp/kubeconfig
```

For debugging purposes, you can increase the log level:

```bash
./bin/manager -zap-log-level=debug --kubeconfig=/tmp/kubeconfig
```

#### Testing with a Sample Resource

Create a sample ShareKube resource to test the operator:

```bash
cat > sample-sharekube.yaml << EOF
apiVersion: sharekube.dev/v1alpha1
kind: ShareKube
metadata:
  name: test-sharekube
  namespace: default
spec:
  targetNamespace: shared-resources
  ttl: "24h"
  resources: []
EOF

kubectl --kubeconfig=/tmp/kubeconfig apply -f sample-sharekube.yaml
```

Verify that the operator has processed the resource:

```bash
# Check if the target namespace was created
kubectl --kubeconfig=/tmp/kubeconfig get namespaces

# Check the status of the ShareKube resource
kubectl --kubeconfig=/tmp/kubeconfig get sharekubes -o yaml
```

You should see the target namespace "shared-resources" created and the ShareKube resource with a status phase of "Ready".

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

The ShareKube project follows a standard Kubernetes operator structure:

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

## Documentation Development

The documentation site is hosted at [docs.sharekube.dev](https://docs.sharekube.dev) and built with Docusaurus.

The documentation includes the following sections:
- [Overview](https://docs.sharekube.dev/overview)
- [Getting Started](https://docs.sharekube.dev/getting-started)
- [API Reference](https://docs.sharekube.dev/api-reference)
- [Future Roadmap](https://docs.sharekube.dev/future-roadmap)
- [Contributing](https://docs.sharekube.dev/contributing)

To run the documentation site locally:

```bash
# Add docs.local.sharekube.dev to your hosts file
# See hosts.txt for the required entry

# Start the documentation site
cd apps/docs
yarn dev

# The site will be available at https://docs.local.sharekube.dev:3000
```

## Release Process

1. **Version Bump**: Update version numbers in the code and manifests
2. **Changelog**: Update the [CHANGELOG.md](CHANGELOG.md) with a summary of changes
3. **Tag Release**: Create a Git tag for the release
4. **CI/CD**: The CI pipeline will build and publish artifacts

## Community

Join our community:

- **GitHub Discussions**: For feature discussions and community help
- **Issues**: Bug reports and feature requests
- **Pull Requests**: Code contributions

## Code of Conduct

We follow the [CNCF Code of Conduct](CODE_OF_CONDUCT.md). Please be respectful and collaborative in all interactions.

## License

ShareKube is licensed under the Apache 2.0 License. All contributions must comply with this license.

---

Thank you for considering a contribution to ShareKube! Your help makes the project better for everyone. 