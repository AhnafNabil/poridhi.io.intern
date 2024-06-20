# Creating Kubernetes Pods: Imperative and Declarative Approaches

In this guide, we will explore how to create Kubernetes Pods using both imperative and declarative methods. This documentation will guide through creating a Pod using each method, highlighting the key steps and differences.

## Imperative Approach

### Overview
The imperative approach involves providing explicit commands or YAML files to Kubernetes to manage objects. This method gives us direct control over the Kubernetes resources.

### Steps to Create a Pod Imperatively

1. **Open a terminal**.
2. **Create a YAML file** called `imperative-pod.yaml` using a text editor (e.g., nano):
   ```sh
   nano imperative-pod.yaml
   ```

3. **Write the Pod specification** in the YAML file:
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: imp-pod
     labels:
       app: myapp
   spec:
     containers:
     - name: imp-container
       image: busybox
       command: ['sh', '-c', 'echo "Welcome to Container MasterClass by CeruleanCanvas" && sleep 60']
   ```

4. **Save and exit** the text editor.

5. **Create the Pod** using the following command:
   ```sh
   kubectl create -f imperative-pod.yaml
   ```

6. **Verify the Pod creation**:
   ```sh
   kubectl get pods
   ```

## Declarative Approach

### Overview
The declarative approach involves providing a configuration file that describes the desired state of the resources. Kubernetes automatically figures out the necessary actions to reach that state.

### Steps to Create a Pod Declaratively

1. **Open another terminal**.
2. **Create a YAML file** called `declarative-pod.yaml` using a text editor:
   ```sh
   nano declarative-pod.yaml
   ```

3. **Write the Pod specification** in the YAML file:
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: decl-pod
     labels:
       app: myapp
   spec:
     containers:
     - name: decl-container
       image: busybox
       command: ['sh', '-c', 'echo "Welcome to Container MasterClass by CeruleanCanvas" && sleep 120']
   ```

4. **Save and exit** the text editor.

5. **Apply the configuration** using the following command:
   ```sh
   kubectl apply -f declarative-pod.yaml
   ```

6. **Verify the Pod creation**:
   ```sh
   kubectl get pods
   ```

## Comparing Imperative and Declarative Pods

### Verify Pod Status
To check the status of both Pods, use the `kubectl get pods` command in both terminals. This command provides a list of Pods with their current status.

### Describe Pods
To get detailed information about each Pod, use the `kubectl describe pods` command followed by the Pod name:
- **Imperative Pod**:
  ```sh
  kubectl describe pods imp-pod
  ```
- **Declarative Pod**:
  ```sh
  kubectl describe pods decl-pod
  ```

### Key Differences
1. **Annotations**:
   - Imperative Pod typically lacks annotations.
   - Declarative Pod may include annotations as it uses a specified Pod template.
   
2. **Container Information**:
   - Both Pods will have unique names and container IDs.
   - They will share the same container image and image ID.
   - Commands executed may differ as specified.

3. **Events Summary**:
   - Kubernetes provides a summary of events for each Pod, including scheduling, image pulling, container creation, and start events.

## Conclusion

By following these steps, we have successfully created and compared Pods using both imperative and declarative methods in Kubernetes. Understanding these approaches helps in managing Kubernetes resources more efficiently, depending on our specific use case and requirements.