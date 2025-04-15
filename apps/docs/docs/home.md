---
id: home
title: ShareKube
slug: /
---

# ShareKube

**ShareKube** is a Kubernetes extension that allows users to create temporary preview environments by copying explicitly defined resources from one namespace to another within the same cluster. This powerful tool is designed for developers and QA teams who need to quickly replicate parts of their environment for testing, demonstration, or validation purposes.

## The Problem ShareKube Solves

**In simple terms:** ShareKube lets you copy parts of your Kubernetes applications to create temporary testing environments. Instead of rebuilding everything from scratch, you can select exactly which pieces you need and ShareKube will create perfect copies that clean themselves up when you're done.

```mermaid
flowchart TD
    subgraph Source["Source Environment"]
        direction TB
        S["Development Environment"]
        S_API["API"]
        S_DB["Database"]
        S_AUTH["Auth"]
        
        S --- S_API
        S --- S_DB
        S --- S_AUTH
    end
    
    Source:::sourceBox
    
    ShareKube[("ShareKube")]
    
    Source --> ShareKube
    
    subgraph Preview["Preview Environments"]
        direction LR
        TeamA["Team A Preview"]
        TeamA_API["API"]
        TeamA --- TeamA_API
        
        TeamB["Team B Preview"]
        TeamB_API["API"]
        TeamB_DB["Database"]
        TeamB --- TeamB_API
        TeamB --- TeamB_DB
        
        QA["QA Environment"]
        QA_API["API"]
        QA_AUTH["Auth"]
        QA --- QA_API
        QA --- QA_AUTH
    end
    
    Preview:::sourceBox
    
    ShareKube --> TeamA
    ShareKube --> TeamB
    ShareKube --> QA
    
    classDef source fill:#3498db,stroke:#333,stroke-width:1px;
    classDef sourceBox fill:none,stroke:#333,stroke-width:1px;
    classDef teamA fill:#2ecc71,stroke:#333,stroke-width:1px;
    classDef teamB fill:#9b59b6,stroke:#333,stroke-width:1px;
    classDef qa fill:#e74c3c,stroke:#333,stroke-width:1px;
    classDef operator fill:white,stroke:#333,stroke-width:2px,color:#222;
    
    class S,S_API,S_DB,S_AUTH source;
    class TeamA,TeamA_API teamA;
    class TeamB,TeamB_API,TeamB_DB teamB;
    class QA,QA_API,QA_AUTH qa;
    class ShareKube operator;
```

### Meet Alex, Platform Lead at FinTech Inc.

Alex leads a team of 12 developers working on a complex microservices-based payment processing platform running on Kubernetes. The platform consists of 30+ services across multiple namespaces:

- `payments-prod`: Contains core payment processing services
- `analytics-prod`: Contains data processing and reporting services
- `auth-prod`: Contains authentication and authorization services

**Alex's Challenge:** The QA team needs an isolated copy of the payment processing services to test a major new feature before it goes to production. They need real services but don't want to replicate the entire cluster.

---

### 🔴 The Challenge: Manual Environment Setup

**Without ShareKube:**

```mermaid
flowchart TD
    A1[Alex needs to<br>create test environment] -->|Without ShareKube| B1[Manual Kubernetes setup]
    B1 -->|Configure YAML files| C1[Copy & modify resources]
    C1 -->|Time-consuming| D1[Resource conflicts]
    C1 -->|Error-prone| E1[Inconsistent environments]
    C1 -->|Costly| F1[Resource waste]
```

Alex would need to:
1. Manually identify all relevant services, deployments, and configs
2. Copy and modify YAML files for each resource
3. Ensure all dependencies are correctly configured
4. Set up network policies and service discovery
5. Manually clean up afterwards

This process would take Alex's team 3-4 days and frequently results in missed dependencies or configuration errors.

---

### 🟢 The Solution: ShareKube Automation

**With ShareKube:**

```mermaid
flowchart TD
    A2[Alex needs to<br>create test environment] -->|With ShareKube| B2[Create ShareKube resource]
    B2 -->|Define resources to copy| C2[ShareKube creates preview]
    C2 -->|Precise control| D2[Only copies what you need]
    C2 -->|Isolated testing| E2[No impact on production]
    C2 -->|Time-to-Live| F2[Automatic cleanup]
```

Alex creates a single ShareKube resource:
```yaml
apiVersion: sharekube.dev/v1alpha1
kind: ShareKube
metadata:
  name: payment-feature-test
  namespace: payments-prod
spec:
  targetNamespace: preview
  ttl: 48h
  resources:
    - kind: Deployment
      name: payment-processor
    - kind: Deployment
      name: transaction-api
    - kind: ConfigMap
      name: payment-config
```

**The Result:** In less than 5 minutes, ShareKube creates an isolated test environment with exactly the components the QA team needs. The environment automatically cleans itself up after 48 hours, and Alex's team saves over 25 hours of DevOps work.

---

## What ShareKube Does

With ShareKube, you can:

- **Create isolated preview environments** by copying selected Kubernetes resources
- **Set time limits** for automatic cleanup of these environments
- **Explicitly define** which resources should be included in your preview
- **Maintain isolation** between your source environment and previews

## How It Works

### 🔄 ShareKube Workflow
The ShareKube workflow is straightforward: create a CRD, specify what to copy, and let the operator handle the rest. When the environment is no longer needed, automatic cleanup occurs based on your TTL settings.

```mermaid
graph LR
    A(Developer creates<br>ShareKube resource) --> C(ShareKube Operator)
    C --> C1(Processes the resource<br>definition)
    C1 --> D(Identifies resources<br>in source namespace)
    D --> E(Creates identical copies<br>in preview namespace)
    E --> F(Monitors and cleans up<br>when TTL expires)
    
    linkStyle 0 stroke:#2ecc71,stroke-width:2px;
    linkStyle 1 stroke:#3498db,stroke-width:2px;
    linkStyle 2 stroke:#e74c3c,stroke-width:2px;
    linkStyle 3 stroke:#9b59b6,stroke-width:2px;
    linkStyle 4 stroke:#f39c12,stroke-width:2px;
```

ShareKube works in five simple steps:

1. **You create a ShareKube resource** that defines what you want to preview
2. **You specify which resources to copy** and how long they should exist
3. **ShareKube locates the resources** in your source namespace
4. **ShareKube creates exact copies** in your preview namespace 
5. **ShareKube automatically cleans up** when the time-to-live expires

With this workflow, you can:
- Create isolated preview environments in seconds
- Test changes safely without affecting production
- Share temporary environments with your team
- Focus on your work while ShareKube handles the infrastructure

## Key Benefits

- **Automation**: Remove manual steps in creating preview environments
- **Isolation**: Test changes without affecting the source environment
- **Efficiency**: Rapidly create and clean up temporary environments
- **Explicit Control**: Only copy the resources you specify

## Future Enhancements

While the current MVP focuses on resource copying, our roadmap includes:

- **Dynamic transformation rules** for automatically modifying copied resources
- **Remote cluster support** for copying resources across clusters
- **Network exposure solutions** for easier access to preview environments

Ready to get started? Check out our [Getting Started guide](getting-started) or dive into the [Architecture Overview](overview). 