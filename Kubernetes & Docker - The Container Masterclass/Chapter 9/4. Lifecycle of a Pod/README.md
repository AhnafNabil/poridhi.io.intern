# Understanding the Lifecycle of Kubernetes Pods

Just like containers, Pods in Kubernetes have distinct lifecycles. This documentation outlines the various states a Pod can go through from creation to termination.

## Pod Lifecycle States

![alt text](Pod-Lifecycle.png)

### 1. Pending
- **Definition**: The Pod has been created and its configurations have been approved by the Kubernetes control plane components (`kube-controller-manager` and `apiserver`).
- **Characteristics**: 
  - The Pod is not yet scheduled to any node.
  - It is waiting for resource allocation and scheduling.

### 2. Running
- **Definition**: The Pod has been scheduled to a node and at least one of its containers is running.
- **Characteristics**:
  - The `kubelet` on the node has accepted the Pod.
  - Containers within the Pod are in the process of running or have already started.

### 3. Succeeded
- **Definition**: All containers in the Pod have successfully completed their tasks and exited gracefully.
- **Characteristics**:
  - The Pod has completed its intended operation.
  - This state is typical for Pods running batch jobs or one-time tasks.

### 4. Failed
- **Definition**: One or more containers in the Pod have terminated unexpectedly due to errors.
- **Characteristics**:
  - This can occur due to reasons like out of memory (OOM) errors or other runtime issues.
  - The Pod transitions to this state from the Running state when failures occur.

### 5. Unknown
- **Definition**: The state of the Pod cannot be determined, usually due to communication issues between the control plane and the node.
- **Characteristics**:
  - The Pod is not running, but the specific reason for this is unclear.
  - This state requires further investigation to diagnose the issue.

## Pod Rescheduling

- **Failed to Pending**: If a Pod is in the Failed state, it can be rescheduled after troubleshooting the issue.
- **Pending to Running**: Upon rescheduling, the Pod returns to the Pending state and proceeds to the Running state once scheduled and started successfully.

## Summary

Understanding the lifecycle of a Pod is crucial for effectively managing and troubleshooting Kubernetes deployments. The lifecycle states of a Pod provide insights into its current status and help identify the actions needed to maintain the desired state of applications running in Kubernetes.

By monitoring the state transitions and handling each state appropriately, we can ensure that your applications are resilient and can recover from failures efficiently.