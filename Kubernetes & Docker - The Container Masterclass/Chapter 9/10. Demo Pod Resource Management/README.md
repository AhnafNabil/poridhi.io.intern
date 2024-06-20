# Managing Pod Resource Requests and Limits in Kubernetes

## Overview

In this guide, we will demonstrate how to specify CPU and memory resources for containers in a Kubernetes Pod. By defining resource requests and limits, you can ensure the Kubernetes scheduler makes informed decisions about Pod placement on nodes and prevent nodes from crashing due to resource exhaustion. We will walk through an example using a Pod with two containers: a MySQL database and a WordPress frontend.

## Prerequisites

- Basic knowledge of Kubernetes concepts.
- `kubectl` command-line tool installed and configured.
- A running Kubernetes cluster.

## Step-by-Step Guide

### 1. Preparing the Pod Specification

We begin by defining a Pod in a YAML file named `resource-pod.yaml`. This Pod includes two containers: one for a MySQL database and another for a WordPress frontend.

#### `resource-pod.yaml`:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: frontend
spec:
  containers:
  - name: mysql
    image: mysql:5.6
    env:
    - name: MYSQL_ROOT_PASSWORD
      value: password
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
  - name: wordpress
    image: wordpress:4.8-apache
    env:
    - name: WORDPRESS_DB_HOST
      value: 127.0.0.1:3306
    - name: WORDPRESS_DB_PASSWORD
      value: password
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
```

### 2. Creating the Pod

Use the `kubectl create` command to create the Pod from the YAML file:

```sh
kubectl create -f resource-pod.yaml
```

### 3. Verifying Pod Status

Check the status of the Pod to ensure it is running:

```sh
kubectl get pods
```

### 4. Troubleshooting Resource Issues

If the Pod is not running correctly, we may need to inspect its status and events for any issues:

```sh
kubectl describe pod frontend
```

In this scenario, the Pod status shows that one of the containers is in a `CrashLoopBackOff` state due to an `OOMKilled` event, indicating it ran out of memory.

### 5. Adjusting Resource Limits

To resolve the issue, we need to increase the memory limits for both containers. Update the `resource-pod.yaml` file as follows:

#### Updated `resource-pod.yaml`:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: frontend
spec:
  containers:
  - name: mysql
    image: mysql:5.6
    env:
    - name: MYSQL_ROOT_PASSWORD
      value: password
    resources:
      requests:
        memory: "256Mi"
        cpu: "250m"
      limits:
        memory: "1Gi"
        cpu: "1"
  - name: wordpress
    image: wordpress:4.8-apache
    env:
    - name: WORDPRESS_DB_HOST
      value: 127.0.0.1:3306
    - name: WORDPRESS_DB_PASSWORD
      value: password
    resources:
      requests:
        memory: "256Mi"
        cpu: "250m"
      limits:
        memory: "1Gi"
        cpu: "1"
```

### 6. Recreating the Pod

First, delete the existing Pod:

```sh
kubectl delete pod frontend
```

Then, create the Pod again with the updated resource limits:

```sh
kubectl create -f resource-pod.yaml
```

### 7. Verifying Pod Status

Check the Pod status again to ensure both containers are running smoothly:

```sh
kubectl get pods
```

### 8. Describing the Pod

Finally, describe the Pod to verify the resource limits have been applied:

```sh
kubectl describe pod frontend
```

We should see that both containers are running without issues, and the events indicate successful creation and running of both containers.

## Conclusion

By setting appropriate resource requests and limits for your containers, we can ensure your applications run smoothly and prevent resource-related issues. This guide demonstrated how to adjust these settings and troubleshoot common issues related to resource constraints.