---
id: overview
title: Architecture Overview
---

# Architecture Overview

ShareKube provides a straightforward yet powerful architecture for creating temporary preview environments in Kubernetes clusters.

## Core Components

### 1. ShareKube Custom Resource Definition (CRD)

The ShareKube CRD is the primary interface for users. It allows you to define:

- **Source resources** to be copied
- **Target namespace** where resources will be deployed
- **Time-to-Live (TTL)** for automatic cleanup
- **Future: Transformation rules** for resource modification

### 2. ShareKube Operator

The ShareKube operator is responsible for:

- Watching for ShareKube CRD instances
- Copying specified resources from source to target namespaces
- Applying any defined transformation rules (future feature)
- Managing the lifecycle of copied resources
- Cleaning up resources when the TTL expires

### 3. Transformation Pipeline (Future Feature)

While not in the current MVP, the transformation pipeline will be responsible for:

- Modifying resources during the copy process
- Removing fields that might cause conflicts (like `spec.clusterIP` for Services)
- Applying custom transformations based on resource type

## Architecture Diagram

```
┌─────────────────────┐      ┌───────────────────────┐
│                     │      │                       │
│  ShareKube CRD      │      │  ShareKube Operator   │
│  (User-defined)     │─────▶│                       │
│                     │      │                       │
└─────────────────────┘      └───────────┬───────────┘
                                         │
                                         ▼
┌─────────────────────┐      ┌───────────────────────┐
│                     │      │                       │
│  Source Namespace   │─────▶│   Target Namespace    │
│  (e.g., "dev")      │      │   (e.g., "preview")   │
│                     │      │                       │
└─────────────────────┘      └───────────────────────┘
```

## Flow of Operation

1. **Define and Apply**: User creates a ShareKube CRD specifying resources to copy
2. **Detection**: The ShareKube operator detects the new CRD
3. **Resource Copying**: The operator locates and copies specified resources
4. **Future: Transformation**: Resources are transformed according to rules
5. **Deployment**: Copied resources are deployed to the target namespace
6. **Lifecycle Management**: The operator tracks the TTL and cleans up resources

## Resource Handling

ShareKube handles various Kubernetes resource types, including:

- Deployments
- Services
- ConfigMaps
- Secrets
- And other standard Kubernetes resources

The operator preserves resource relationships and dependencies during the copying process to ensure proper functionality in the target namespace.

Each resource is copied with minimal modifications, preserving the original configuration as much as possible while ensuring it can function correctly in the new namespace. 