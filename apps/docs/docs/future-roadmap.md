---
id: future-roadmap
title: Future Roadmap
---

# Future Roadmap

ShareKube's current MVP focuses on the core functionality of copying resources between namespaces. We have an ambitious roadmap of features planned for future releases to enhance its capabilities.

## Dynamic Transformation Rules

A key upcoming feature is the ability to dynamically transform resources during the copy process. This will include:

- **Field Removal**: Automatically remove fields that might cause conflicts (e.g., `spec.clusterIP` for Services)
- **Field Modification**: Transform values in the copied resources (e.g., change labels or annotations)
- **Prefix/Suffix Addition**: Add prefixes or suffixes to resource names to avoid conflicts
- **Plugin Architecture**: Allow custom transformations through a plugin system

Example future transformation rules:

```yaml
spec:
  transformationRules:
    # Remove specific fields from Service resources
    - kind: Service
      removeFields:
        - spec.clusterIP
        - spec.clusterIPs
        - status
    
    # Add a prefix to all Deployment names
    - kind: Deployment
      nameTransform:
        prefix: "preview-"
    
    # Replace environment variables in ConfigMaps
    - kind: ConfigMap
      dataTransforms:
        - field: "data.env-vars"
          replace:
            - pattern: "PROD_"
              with: "PREVIEW_"
```

## Multi-Cluster Support

We plan to extend ShareKube to work across multiple Kubernetes clusters:

- **Remote Cluster Targeting**: Copy resources from one cluster to another
- **Credential Management**: Securely store and use credentials for remote clusters
- **Cluster Discovery**: Automatically discover available target clusters

Example configuration:

```yaml
spec:
  targetCluster:
    name: staging-cluster
    kubeconfigSecret: staging-kubeconfig
  targetNamespace: preview
  # ... other fields
```

## Network Exposure Integration

To make preview environments more accessible, we'll add:

- **Automatic Ingress Generation**: Create Ingress rules for exposed services
- **Dynamic DNS Configuration**: Set up DNS entries for preview environments
- **TLS Certificate Provisioning**: Automatic TLS certificate generation for secure access

## Enhanced Lifecycle Management

Improved management of preview environment lifecycles:

- **Graduated TTL**: Different TTLs for different resources
- **Usage-Based Expiry**: Expire previews based on activity rather than just time
- **Scheduled Creation**: Create preview environments on a schedule
- **Event Triggers**: Create or update previews based on external events (CI/CD, Git hooks)

## User Interface and Dashboard

A web-based interface for managing ShareKube:

- **Visual Preview Creation**: Create previews through a UI instead of YAML
- **Status Dashboard**: View all active previews and their status
- **Resource Utilization**: Monitor resource usage of preview environments
- **Access Management**: Control who can create and access previews

## Integration with Development Workflows

- **CI/CD Integration**: Create previews as part of CI/CD pipelines
- **Pull Request Previews**: Automatically create previews for pull requests
- **Test Environment Management**: Integration with testing frameworks

## Observability and Monitoring

- **Enhanced Metrics**: Detailed metrics on preview usage and performance
- **Audit Logging**: Track creation, access, and deletion of previews
- **Cost Tracking**: Monitor cloud costs associated with previews

## Feedback and Contributions

We welcome feedback on our roadmap and contributions to help implement these features. If you have ideas for additional features or would like to contribute to development, please see our [Contributing](contributing) guide. 