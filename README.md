# ShareKube

<div align="center">
  <img src="packages/shared-assets/logo/logo.svg" alt="ShareKube Logo" width="200"/>
  
  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
  [![Documentation](https://img.shields.io/badge/docs-docs.sharekube.dev-brightgreen)](https://docs.sharekube.dev)
</div>

## Overview

ShareKube is a Kubernetes extension that allows users to create temporary preview environments by copying explicitly defined resources from one namespace to another within the same cluster. Perfect for development teams, CI/CD pipelines, and collaborative Kubernetes workflows.

## 🔄 How ShareKube Works

ShareKube follows a simple workflow that makes sharing Kubernetes resources between namespaces easy:

1. **Resource Selection**: Define which resources you want to share from your source namespace
2. **Create ShareKube Object**: Create a ShareKube custom resource specifying target namespace and TTL
3. **Automatic Copying**: ShareKube operator copies the selected resources to the target namespace
4. **Preview Environment Ready**: Use the copied resources in the target namespace for testing/review
5. **Automatic Cleanup**: After the TTL expires, all copied resources are automatically removed

This workflow enables powerful use cases:
- Creating per-PR preview environments in CI/CD pipelines
- Sharing development environments with team members
- Setting up temporary demo environments for stakeholders
- Creating isolated testing environments that mirror production

## ✨ Key Features

- **Quick Preview Environments**: Create temporary namespaces with copies of your resources
- **Selective Resource Copying**: Choose exactly which resources to share from source namespaces
- **Time-to-Live (TTL)**: Automatically clean up resources after a configurable period
- **Non-Destructive**: Original resources remain untouched in their source namespaces
- **Kubernetes Native**: Uses standard CRDs and follows Kubernetes patterns

## 🚀 Quick Start

```bash
# Install the ShareKube operator
kubectl apply -f https://github.com/miloszsobczak/sharekube/releases/latest/download/sharekube-operator.yaml

# Create a ShareKube resource to share resources
cat <<EOF | kubectl apply -f -
apiVersion: sharekube.dev/v1alpha1
kind: ShareKube
metadata:
  name: demo-preview
  namespace: default
spec:
  targetNamespace: preview-env
  ttl: "24h"
  resources:
    - kind: Deployment
      name: my-application
    - kind: Service
      name: my-application-svc
EOF

# Check that resources were copied
kubectl get deployments,services -n preview-env
```

## 📚 Documentation

For comprehensive documentation, visit [docs.sharekube.dev](https://docs.sharekube.dev)

- [Overview](https://docs.sharekube.dev/overview)
- [Getting Started](https://docs.sharekube.dev/getting-started)
- [API Reference](https://docs.sharekube.dev/api-reference)
- [Future Roadmap](https://docs.sharekube.dev/future-roadmap)
- [Contributing](https://docs.sharekube.dev/contributing)

## 🔧 Installation

```bash
# Install the ShareKube operator
kubectl apply -f https://github.com/miloszsobczak/sharekube/releases/latest/download/sharekube-operator.yaml
```

## 💻 Development

### Prerequisites

- Node.js 16+
- Yarn
- Go 1.19+ (for operator development)
- Kubernetes cluster for testing
- [mkcert](https://github.com/FiloSottile/mkcert) (for local HTTPS development)

### Local Setup

```bash
# Clone the repository
git clone https://github.com/miloszsobczak/sharekube.git
cd sharekube

# Install dependencies
yarn

# Set up local HTTPS certificates (for documentation site)
./scripts/setup-local-https.sh
```

For detailed development instructions, see [CONTRIBUTING.md](CONTRIBUTING.md).

## 🌟 Project Structure

This is a monorepo containing all ShareKube projects:

- `apps/` - Applications
  - `docs/` - Documentation site (Docusaurus)
- `packages/` - Shared packages
  - `operator/` - Kubernetes operator
  - `shared-assets/` - Shared assets like logos and icons

## 📝 Project Documentation

- [CONTRIBUTING.md](CONTRIBUTING.md): Contribution guidelines and development setup
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md): Our community code of conduct
- [SECURITY.md](SECURITY.md): Security policy and vulnerability reporting
- [CHANGELOG.md](CHANGELOG.md): Version history and changes
- [LICENSE](LICENSE): Apache 2.0 license details

## 👥 Community

- [GitHub Issues](https://github.com/miloszsobczak/sharekube/issues): Bug reports & feature requests
- [GitHub Discussions](https://github.com/miloszsobczak/sharekube/discussions): Questions & community discussions

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to submit contributions.

All contributors are expected to follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## 📜 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details. 