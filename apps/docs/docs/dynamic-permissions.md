---
id: dynamic-permissions
title: Dynamic Permissions
---

# Dynamic Permissions in ShareKube

ShareKube now supports dynamic RBAC permissions to enhance security when copying resources between namespaces. This feature allows for fine-grained control over what permissions are granted, reducing the potential attack surface.

## How It Works

Instead of requiring broad cluster-wide permissions, ShareKube dynamically creates and manages Role and RoleBinding resources specific to each ShareKube instance, granting only the permissions needed for the requested resource types.

The permissions are:

1. **Source namespace**: Read-only permissions for the specific resource types being copied
2. **Target namespace**: Full permissions for the specific resource types being created

## Enabling Dynamic Permissions

Dynamic permissions can be enabled for a ShareKube resource by adding the `accessControl` field to your ShareKube spec:

```yaml
apiVersion: sharekube.dev/v1alpha1
kind: ShareKube
metadata:
  name: sample-app-preview
  namespace: dev
spec:
  targetNamespace: preview-app-123
  ttl: "24h"
  resources:
    - kind: Deployment
      name: sample-api
    - kind: Service
      name: sample-api-svc
  accessControl:
    restrict: true
    allowedSourceNamespaces:
      - dev
      - staging
    allowedTargetNamespaces:
      - preview-*
```

## Configuration Options

The `accessControl` section supports the following options:

- `restrict`: When set to `true`, enables dynamic permission creation based on resource types
- `allowedSourceNamespaces`: List of namespaces that can be used as sources (currently informational only)
- `allowedTargetNamespaces`: List of namespaces that can be used as targets (currently informational only)

## Benefits

- **Enhanced Security**: Only grants the minimum necessary permissions
- **Better Compliance**: Easier to meet security compliance requirements
- **Reduced Attack Surface**: Limits potential damage from security breaches
- **Namespace Isolation**: Maintains proper namespace boundaries

## Implementation Details

The dynamic permissions are created as follows:

1. When a ShareKube resource is created, the controller analyzes the requested resources
2. It creates appropriate Role and RoleBinding resources with permissions for only those resource types
3. Permissions are tracked in the ShareKube status via the `dynamicPermissions` field
4. When the ShareKube resource is deleted, these dynamic permissions are automatically cleaned up

## Operator Permissions

The ShareKube operator now requires only the following cluster-wide permissions:

- Managing ShareKube CRDs
- Managing namespaces (creation only)
- Managing Roles and RoleBindings (to create the dynamic permissions)
- Leader election leases

All other permissions are granted dynamically at the namespace level.

## Future Enhancements

Future versions of ShareKube will implement:

- Validation of allowedSourceNamespaces and allowedTargetNamespaces
- RBAC auditing and policy enforcement
- Integration with external policy engines
- Fine-grained resource field restrictions 