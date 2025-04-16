# ShareKube Basic Example

This directory contains a basic example to test ShareKube functionality.

## Prerequisites

- A Kubernetes cluster
- ShareKube operator installed
- kubectl configured to access your cluster

## Getting Started

1. Create the required namespaces:

```bash
kubectl create namespace dev
kubectl create namespace preview
```

2. Deploy the sample application components to the dev namespace:

```bash
kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
```

3. Verify the resources exist in the dev namespace:

```bash
kubectl get deployment -n dev
kubectl get configmap -n dev
```

4. Create the ShareKube resource to clone the application to the preview namespace:

```bash
kubectl apply -f sharekube.yaml
```

5. Verify the resources were copied to the preview namespace:

```bash
kubectl get deployment -n preview
kubectl get configmap -n preview
```

6. Check the ShareKube resource status:

```bash
kubectl get sharekube -n dev sample-app-preview -o yaml
```

Look for the `status.phase` field, which should be set to `Ready` when the preview environment has been successfully created.

## Testing Changes

Try modifying resources in the preview namespace without affecting the original development environment:

```bash
# Scale the frontend deployment
kubectl scale deployment sample-frontend -n preview --replicas=2

# Update the ConfigMap
kubectl patch configmap sample-app-config -n preview --type merge -p '{"data":{"app.settings":"{\"apiEndpoint\":\"/api/v2\",\"logLevel\":\"debug\",\"enableCache\":\"true\"}"}}'
```

## Cleanup

To clean up the resources:

```bash
# Delete the ShareKube resource
kubectl delete -f sharekube.yaml

# Delete the namespaces
kubectl delete namespace dev
kubectl delete namespace preview
``` 