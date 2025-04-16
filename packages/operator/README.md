# ShareKube Operator

ShareKube is a Kubernetes operator that allows creating temporary preview environments by copying resources from one namespace to another.

## Overview

The ShareKube operator watches for ShareKube custom resources, which specify a set of Kubernetes resources to copy from a source namespace to a target namespace. This is useful for creating temporary environments for testing, development, or demonstration purposes.

## Features

- **Resource Copying**: Copy specific Kubernetes resources from one namespace to another
- **TTL-based Cleanup**: Set a Time-to-Live (TTL) for automatic cleanup of preview environments
- **Explicit Resource Control**: Specify exactly which resources should be copied
- **Resource Tracking**: Resources are labeled to track ownership for proper cleanup
- **Resource Transformation**: Automatic handling of cluster-specific fields (like Service ClusterIPs)
- **Future Features**:
  - Advanced transformation rules for resources during copying
  - Remote cluster support for copying between clusters

## Supported Resources

ShareKube supports copying the following resource types:
- Deployments
- Services
- ConfigMaps
- Secrets
- Other Kubernetes resources via the dynamic client

## Installation

### Prerequisites

- Kubernetes cluster (v1.16+)
- kubectl configured to access your cluster
- Cluster admin permissions (to install CRDs)

### Install with kubectl

1. Install the CRD:

```bash
kubectl apply -f config/crd/bases/sharekube.dev_sharekubes.yaml
```

2. Install the operator:

```bash
kubectl apply -f config/manager/manager.yaml
```

## Usage

1. Create a ShareKube resource:

```yaml
apiVersion: sharekube.dev/v1alpha1
kind: ShareKube
metadata:
  name: my-preview
  namespace: dev
spec:
  targetNamespace: preview  # Where resources will be copied to
  ttl: 1h                   # How long the resources will exist
  resources:
    # Copy deployment from same namespace as this ShareKube resource
    - kind: Deployment
      name: my-app
    
    # Copy service from same namespace
    - kind: Service
      name: my-app-svc
    
    # Copy ConfigMap from a specific namespace 
    - kind: ConfigMap
      name: my-app-config
      namespace: config
```

2. The operator will create the target namespace if it doesn't exist
3. Resources will be copied to the target namespace with appropriate labels
4. After the TTL expires, the operator will automatically clean up the copied resources

## Resource Handling

When a resource is copied to the target namespace:

1. A copy of the original resource is created in the target namespace
2. The copied resource is labeled with `sharekube.dev/copied: "true"` for tracking
3. Cluster-specific fields (like ClusterIPs for Services) are cleared to avoid conflicts
4. Resource metadata is updated to match the target namespace

## Development

### Building

To build the operator:

```bash
# Build the binary
go build -o bin/manager main.go

# Build the Docker image
docker build -t sharekube/operator:dev .
```

### Running Locally

To run the operator locally:

```bash
# Install CRDs
kubectl apply -f config/crd/bases/sharekube.dev_sharekubes.yaml

# Run the operator
go run main.go
```

### Project Structure

```
├── api/                  # API definitions (CRDs)
│   └── v1alpha1/         # API version
├── config/               # Operator configuration
│   ├── crd/              # CRD manifests
│   ├── manager/          # Controller manager deployment
│   └── samples/          # Sample CR files
├── controllers/          # Operator controllers
│   └── sharekube_controller.go  # Main reconciliation logic
├── pkg/                  # Shared packages
│   └── resources/        # Resource management
│       └── handler.go    # Resource copying implementation
```

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.