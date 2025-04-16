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
- **Future: Advanced transformation rules** for resource modification

### 2. ShareKube Operator

The ShareKube operator is responsible for:

- Watching for ShareKube CRD instances
- Copying specified resources from source to target namespaces
- Applying basic transformations for compatibility (e.g., clearing Service ClusterIPs)
- Managing the lifecycle of copied resources
- Cleaning up resources when the TTL expires

### 3. Transformation Pipeline

The transformation pipeline handles modifications to resources during the copy process:

**Current Features:**
- Automatic removal of cluster-specific fields (like `spec.clusterIP` for Services)
- Updating namespace references
- Adding tracking labels for lifecycle management

**Future Features:**
- User-defined transformation rules via the CRD
- Custom field removal based on resource type
- Advanced transformations based on resource relationships

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
4. **Transformation**: Resources are automatically transformed for compatibility
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