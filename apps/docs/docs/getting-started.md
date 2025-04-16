---
id: getting-started
title: Getting Started
---

# Getting Started with ShareKube

ShareKube is a Kubernetes operator that enables you to create isolated preview environments by copying resources from one namespace to another. This guide will help you install ShareKube on your Kubernetes cluster and create your first preview environment.

## Prerequisites

Before you begin, make sure you have:

- A Kubernetes cluster (v1.16+)
- `kubectl` installed and configured to access your cluster
- Cluster admin permissions (to install CRDs)

## Installation

### Install with kubectl

1. Install the ShareKube CRD:

```bash
kubectl apply -f https://raw.githubusercontent.com/miloszsobczak/sharekube/main/packages/operator/config/crd/bases/sharekube.dev_sharekubes.yaml
```

2. Install the ShareKube operator:

```bash
kubectl apply -f https://raw.githubusercontent.com/miloszsobczak/sharekube/main/packages/operator/config/manager/manager.yaml
```

The operator will be installed in the `sharekube-system` namespace by default.

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

## Ready-to-Use Examples

For your convenience, we've provided ready-to-use examples to help you test ShareKube quickly:

1. **Basic Example**: A minimal setup with a simple deployment and ConfigMap
2. **Microservice Example**: A complete multi-component application with frontend, backend, and database

These examples are located in the [samples directory](https://github.com/miloszsobczak/sharekube/tree/main/samples) of the repository, and each includes detailed instructions and all the YAML files needed.

To use these examples, clone the repository and follow the instructions in the sample directories:

```bash
git clone https://github.com/miloszsobczak/sharekube.git
cd sharekube/samples
```

## Your First Preview Environment

If you prefer to build your own example step by step, let's create a practical example that demonstrates how ShareKube can be used to create preview environments for a microservice application.

### 1. Create a Sample Microservice Application

We'll create a simple microservice application with three components:
- A frontend service
- A backend API service
- A database

First, let's create the necessary namespaces:

```bash
# Create the development namespace where our original application will run
kubectl create namespace dev

# Create the preview namespace where the preview environment will be created
kubectl create namespace preview
```

Now, let's deploy our sample microservice application to the `dev` namespace:

<details>
<summary>Click to expand sample application YAML definitions</summary>

```bash
# Deploy the database
kubectl -n dev apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-db
  namespace: dev
  labels:
    app: sample-app
    component: database
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sample-app
      component: database
  template:
    metadata:
      labels:
        app: sample-app
        component: database
    spec:
      containers:
      - name: postgres
        image: postgres:14
        env:
        - name: POSTGRES_PASSWORD
          value: "password123"
        - name: POSTGRES_USER
          value: "sampleuser"
        - name: POSTGRES_DB
          value: "sampledb"
        ports:
        - containerPort: 5432
---
apiVersion: v1
kind: Service
metadata:
  name: sample-db-svc
  namespace: dev
  labels:
    app: sample-app
    component: database
spec:
  selector:
    app: sample-app
    component: database
  ports:
  - port: 5432
    targetPort: 5432
EOF

# Deploy the backend API
kubectl -n dev apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-api
  namespace: dev
  labels:
    app: sample-app
    component: api
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sample-app
      component: api
  template:
    metadata:
      labels:
        app: sample-app
        component: api
    spec:
      containers:
      - name: api
        image: nginx:latest  # Replace with your actual API image
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "sample-db-svc"
        - name: DB_PORT
          value: "5432"
        - name: DB_USER
          value: "sampleuser"
        - name: DB_PASSWORD
          value: "password123"
        - name: DB_NAME
          value: "sampledb"
---
apiVersion: v1
kind: Service
metadata:
  name: sample-api-svc
  namespace: dev
  labels:
    app: sample-app
    component: api
spec:
  selector:
    app: sample-app
    component: api
  ports:
  - port: 8080
    targetPort: 8080
EOF

# Deploy the frontend
kubectl -n dev apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-frontend
  namespace: dev
  labels:
    app: sample-app
    component: frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sample-app
      component: frontend
  template:
    metadata:
      labels:
        app: sample-app
        component: frontend
    spec:
      containers:
      - name: frontend
        image: nginx:latest  # Replace with your actual frontend image
        ports:
        - containerPort: 80
        env:
        - name: API_URL
          value: "http://sample-api-svc:8080"
---
apiVersion: v1
kind: Service
metadata:
  name: sample-frontend-svc
  namespace: dev
  labels:
    app: sample-app
    component: frontend
spec:
  selector:
    app: sample-app
    component: frontend
  ports:
  - port: 80
    targetPort: 80
EOF

# Create a ConfigMap with application settings
kubectl -n dev apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: sample-app-config
  namespace: dev
  labels:
    app: sample-app
data:
  app.settings: |
    {
      "apiEndpoint": "/api/v1",
      "logLevel": "info",
      "enableCache": "true"
    }
EOF
```

</details>

Verify that all components are running:

```bash
kubectl get all -n dev -l app=sample-app
```

### 2. Create a ShareKube Resource

Now, let's create a ShareKube resource to clone our application to the preview namespace. The ShareKube operator will copy the specified resources from the source namespace (`dev`) to the target namespace (`preview`):

```bash
kubectl apply -f - <<EOF
apiVersion: sharekube.dev/v1alpha1
kind: ShareKube
metadata:
  name: sample-app-preview
  namespace: dev
spec:
  targetNamespace: preview
  ttl: 2h
  resources:
    # Database
    - kind: Deployment
      name: sample-db
    - kind: Service
      name: sample-db-svc
    
    # Backend API
    - kind: Deployment
      name: sample-api
    - kind: Service
      name: sample-api-svc
    
    # Frontend
    - kind: Deployment
      name: sample-frontend
    - kind: Service
      name: sample-frontend-svc
    
    # Configuration
    - kind: ConfigMap
      name: sample-app-config
EOF
```

### 3. Verify Your Preview Environment

Check if the ShareKube resource was processed:

```bash
kubectl get sharekube sample-app-preview -n dev -o yaml
```

Look for the `status.phase` field, which should be set to `Ready` when the preview environment has been successfully created.

Verify that all resources were copied to the preview namespace:

```bash
# List all preview namespace resources
kubectl get all -n preview -l app=sample-app

# Check the ConfigMap
kubectl get configmap -n preview -l app=sample-app
```

### 4. Access Your Preview Environment

:::caution Coming Soon
Additional options for accessing preview environments (including automatic Ingress creation, dynamic DNS configuration, and service mesh integration) will be added in upcoming releases. Stay tuned!
:::

In a real-world scenario, you would typically create an Ingress or use a service mesh to expose your application. For demonstration purposes, you can use port-forwarding to access the frontend:

```bash
# Forward the frontend service port
kubectl port-forward -n preview svc/sample-frontend-svc 8080:80
```

Now you can access your preview environment at http://localhost:8080.

### 5. Testing Changes in the Preview Environment

You can make changes directly to resources in the preview namespace to test modifications without affecting the original development environment. For example:

```bash
# Scale the frontend deployment
kubectl scale deployment sample-frontend -n preview --replicas=2

# Update the ConfigMap
kubectl patch configmap sample-app-config -n preview --type merge -p '{"data":{"app.settings":"{\"apiEndpoint\":\"/api/v2\",\"logLevel\":\"debug\",\"enableCache\":\"true\"}"}}'
```

These changes will only affect the preview environment and won't impact the original application in the `dev` namespace.

## Cleanup

The preview environment will be automatically deleted after the TTL (2 hours in this example) expires. If you want to delete it manually, simply delete the ShareKube resource:

```bash
kubectl delete sharekube sample-app-preview -n dev
```

You can also directly delete the preview namespace for a complete cleanup:

```bash
kubectl delete namespace preview
```

## Next Steps

Now that you've set up ShareKube and created your first preview environment, you can:

- Learn more about the [ShareKube architecture](overview)
- Explore the [CRD API Reference](api-reference)
- Check out the [Future Roadmap](future-roadmap)
- [Contribute](contributing) to the ShareKube project - See our [Contributing Guide](contributing) for information on setting up a development environment 