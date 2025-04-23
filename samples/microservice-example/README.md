# ShareKube Microservice Example

This example demonstrates how to use ShareKube with a complete microservice application consisting of:

- Frontend service
- Backend API service
- Database (PostgreSQL)
- Application configuration

## Prerequisites

- A Kubernetes cluster
- ShareKube operator installed
- kubectl configured to access your cluster
- kustomize (optional, but recommended)

## Getting Started

### Using Kustomize (Recommended)

1. Create the required namespaces:

```bash
kubectl create namespace dev
kubectl create namespace preview
```

Alternatively, uncomment the namespace.yaml in kustomization.yaml and apply everything together.

2. Deploy the entire application at once:

```bash
kubectl apply -k .
```

### Without Kustomize

1. Create the required namespaces:

```bash
kubectl apply -f namespace.yaml
```

2. Deploy the individual components:

```bash
kubectl apply -f database.yaml
kubectl apply -f backend.yaml
kubectl apply -f frontend.yaml
kubectl apply -f configmap.yaml
```

3. Create the ShareKube resource:

```bash
kubectl apply -f sharekube.yaml
```

## Verifying Your Setup

1. Check that all application components are running in the dev namespace:

```bash
kubectl get all -n dev -l app=sample-app
```

2. Check the ShareKube resource status:

```bash
kubectl get sharekube -n dev sample-app-preview -o yaml
```

Look for the `status.phase` field, which should be set to `Ready` when the preview environment has been successfully created.

3. Verify that all resources were copied to the preview namespace:

```bash
kubectl get all -n preview -l app=sample-app
kubectl get configmap -n preview -l app=sample-app
```

4. Check the dynamic RBAC permissions that were created:

```bash
# Check the dynamic permissions in the status
kubectl get sharekube -n dev sample-app-preview -o=jsonpath='{.status.dynamicPermissions}'

# List the dynamic roles
kubectl get roles -n dev | grep sample-app-preview
kubectl get roles -n preview | grep sample-app-preview

# Examine role details and permissions
kubectl describe role sharekube-sample-app-preview-source -n dev
kubectl describe role sharekube-sample-app-preview-target -n preview
```

The dynamic permissions ensure that the operator only has the minimum necessary permissions to copy and manage the specific resources you've defined in your ShareKube resource.

## Accessing the Application

You can use port forwarding to access the frontend service:

```bash
# Access the original app
kubectl port-forward -n dev svc/sample-frontend-svc 8080:80

# Access the preview environment
kubectl port-forward -n preview svc/sample-frontend-svc 8081:80
```

Now you can access:
- Original app at: http://localhost:8080
- Preview environment at: http://localhost:8081

## Testing Changes

You can make changes to resources in the preview namespace without affecting the original environment:

```bash
# Scale the frontend deployment
kubectl scale deployment sample-frontend -n preview --replicas=2

# Update the backend API environment variable
kubectl set env deployment/sample-api -n preview API_VERSION=v2

# Update the ConfigMap
kubectl patch configmap sample-app-config -n preview --type merge -p '{"data":{"app.settings":"{\"apiEndpoint\":\"/api/v2\",\"logLevel\":\"debug\",\"enableCache\":\"true\"}"}}'
```

## Cleanup

To clean up all resources:

```bash
# Using kustomize
kubectl delete -k .

# Or manually
kubectl delete -f sharekube.yaml
kubectl delete -f frontend.yaml
kubectl delete -f backend.yaml
kubectl delete -f database.yaml
kubectl delete -f configmap.yaml
kubectl delete -f namespace.yaml
``` 