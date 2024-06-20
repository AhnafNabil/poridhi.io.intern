# Kubernetes Deployment and Rolling Update Operations

## Introduction

**Rolling Update**:
A rolling update is a deployment strategy in Kubernetes where updates are applied incrementally to the Pods in a deployment, ensuring that the application remains available throughout the update process. This strategy helps minimize downtime and ensures seamless updates.

![rolling-updates](./images/rolling-definition.jpg)

**Rollback**:
Rollback is the process of reverting a deployment to a previous state or version. In Kubernetes, rollback is used to undo changes made during an update and restore the application to a known working state. This is useful in case of errors or issues introduced by an update, allowing operators to quickly restore service reliability.

![rollback](./images/rollback-definition.jpg)

## Steps

### Step 1: Create an nginx Deployment

1. Open the `update-pod.yaml` file, which defines an nginx deployment with 10 replicas using nginx 1.7.9:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy-nginx
spec:
  replicas: 10
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

2. Create the deployment:

```sh
kubectl apply -f update-pod.yaml
```

3. Verify that the deployment is created and all Pods are up-to-date and available:

```sh
kubectl get pods
```

### Step 2: Describe the Deployment

1. Describe the deployment to check its configuration and ensure that it was created successfully:

```sh
kubectl describe deployment deploy-nginx
```

### Step 3: Update the Deployment Image

1. Update the deployment to use the nginx 1.9.1 image:

```sh
kubectl set image deployment/deploy-nginx nginx=nginx:1.9.1
```

2. Verify the update status:

```sh
kubectl rollout status deployment/deploy-nginx
```

- Initially, some of the replicas will be updated. Wait until all replicas are updated.

### Step 4: Perform Another Update

1. Update the deployment to use the `nginx:alpine` image:

```sh
kubectl set image deployment/deploy-nginx nginx=nginx:alpine
```

### Step 5: Check Rollout History

1. Get the rollout history of the deployment:

```sh
kubectl rollout history deployment/deploy-nginx
```

2. Check the details of a specific revision (e.g., revision 2):

```sh
kubectl rollout history deployment/deploy-nginx --revision=2
```

- This command shows that revision 2 involved updating the Pod template with the image `nginx:1.9.1`.

### Step 6: Undo the Update

1. Undo the last update to the deployment:

```sh
kubectl rollout undo deployment/deploy-nginx
```

2. Verify that all 10 Pods are up and running:

```sh
kubectl get pods
```

### Step 7: Roll Back to a Specific Revision

1. Roll back the deployment to revision 2:

```sh
kubectl rollout undo deployment/deploy-nginx --to-revision=2
```

2. Describe the deployment to verify that the image is set to `nginx:1.9.1`:

```sh
kubectl describe deployment deploy-nginx
```

3. Check the events to see the scaling actions:

```sh
kubectl get events --sort-by=.metadata.creationTimestamp
```

- The events will show the deployment being scaled up and down multiple times, which is expected as part of the rolling update and rollback processes.