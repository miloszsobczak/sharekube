---
id: getting-started
title: Getting Started
---

# Getting Started with ShareKube

This guide will help you install and set up ShareKube on your Kubernetes cluster.

## Prerequisites

Before you begin, make sure you have:

- A Kubernetes cluster (v1.16+)
- `kubectl` installed and configured to access your cluster
- Cluster admin permissions (to install CRDs)

## Installation

### Option 1: Install with kubectl

1. Install the ShareKube CRD:

```bash
kubectl apply -f https://raw.githubusercontent.com/miloszsobczak/sharekube/main/deploy/crds/sharekube-crd.yaml
```

2. Install the ShareKube operator:

```bash
kubectl apply -f https://raw.githubusercontent.com/miloszsobczak/sharekube/main/deploy/operator.yaml
```

### Option 2: Install with Helm

1. Add the ShareKube Helm repository:

```bash
helm repo add sharekube https://miloszsobczak.github.io/sharekube/charts
helm repo update
```

2. Install ShareKube:

```bash
helm install sharekube sharekube/sharekube --namespace sharekube --create-namespace
```

## Verify Installation

Check that the ShareKube operator is running:

```bash
kubectl get pods -n sharekube-system
```

You should see output similar to:

```
NAME                                 READY   STATUS    RESTARTS   AGE
sharekube-operator-f84976b4d-x9pqr   1/1     Running   0          30s
```

## Your First Preview Environment

### 1. Create Source Resources

For this example, let's create a simple deployment in the `dev` namespace:

```bash
# Create the dev namespace if it doesn't exist
kubectl create namespace dev

# Create a simple deployment
kubectl -n dev apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-app
  namespace: dev
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hello
  template:
    metadata:
      labels:
        app: hello
    spec:
      containers:
      - name: hello
        image: nginx:latest
        ports:
        - containerPort: 80
EOF

# Create a service for the deployment
kubectl -n dev apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: hello-svc
  namespace: dev
spec:
  selector:
    app: hello
  ports:
  - port: 80
    targetPort: 80
EOF
```

### 2. Create the Target Namespace

```bash
kubectl create namespace preview
```

### 3. Create a ShareKube CRD

Now, let's create a ShareKube CRD to copy the resources to the `preview` namespace:

```bash
kubectl apply -f - <<EOF
apiVersion: sharekube.dev/v1alpha1
kind: ShareKube
metadata:
  name: my-preview
  namespace: dev
spec:
  targetNamespace: preview
  ttl: 1h
  resources:
    - kind: Deployment
      name: hello-app
    - kind: Service
      name: hello-svc
EOF
```

### 4. Verify the Preview Environment

Check if the resources were copied:

```bash
# Check deployments
kubectl get deployments -n preview

# Check services
kubectl get services -n preview
```

You should see the `hello-app` deployment and `hello-svc` service in the `preview` namespace.

## Cleanup

The preview environment will be automatically deleted after the TTL (1 hour in this example) expires. If you want to delete it manually, simply delete the ShareKube CRD:

```bash
kubectl delete sharekube my-preview -n dev
```

## Next Steps

Now that you've set up ShareKube and created your first preview environment, you can:

- Learn more about the [ShareKube architecture](overview)
- Explore the [CRD API Reference](api-reference)
- Check out the [Future Roadmap](future-roadmap)
- [Contribute](contributing) to the ShareKube project 