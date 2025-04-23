# ShareKube Examples

This directory contains examples that demonstrate how to use ShareKube in different scenarios. These examples can be used both for testing and learning purposes.

## Available Examples

### Basic Example

A minimal example to test ShareKube's core functionality. It includes:
- A simple Deployment
- A ConfigMap
- A ShareKube resource
- Dynamic RBAC permissions

[Go to Basic Example →](basic-example/)

### Microservice Example

A complete microservice application to demonstrate ShareKube in a realistic scenario. It includes:
- A frontend service
- A backend API service
- A database
- Application configuration
- Dynamic RBAC permissions
- Kustomization for easy deployment

[Go to Microservice Example →](microservice-example/)

### Dynamic Permissions Example

A focused example demonstrating ShareKube's dynamic RBAC permissions system. It shows:
- Fine-grained permission creation
- Namespace-specific role bindings
- Restricted access to resources

[View Example →](sharekube-with-dynamic-permissions.yaml)

## Using the Examples

Each example directory contains:
- Kubernetes resource YAML files
- A README with detailed instructions
- Easy-to-follow steps to set up, test, and clean up

## Prerequisites

- A Kubernetes cluster
- ShareKube operator installed
- kubectl configured to access your cluster

For more information on getting started with ShareKube, see the [Getting Started Guide](../apps/docs/docs/getting-started.md). 