Filtering output with custom columns in `kubectl get` allows you to specify exactly which fields you want to see from the resource objects fetched from the Kubernetes API. This is particularly useful when you only need specific information and want to customize the output format. Here's a step-by-step guide on how to use the `custom-columns` flag:

### Procedure to Filter Output with Custom Columns:

1. **Basic `kubectl get` Command:**
   Start with a basic `kubectl get` command to retrieve resource objects from Kubernetes. For example, to list pods:
   ```bash
   kubectl get pods
   ```

2. **Adding `custom-columns` Flag:**
   Use the `-o custom-columns=` flag followed by a comma-separated list of column specifiers in the format `<header-name>:<json-path-expression>`. This allows you to define custom headers and specify which fields (using JSONPath expressions) to display.

   Syntax:
   ```bash
   kubectl get <resource> -o custom-columns=<header-name>:<json-path-expression>,<header-name>:<json-path-expression>,...
   ```

   Replace `<resource>` with the Kubernetes resource you are querying (e.g., `pods`, `services`, `deployments`).

3. **JSONPath Expressions:**
   JSONPath expressions are used to extract specific fields from Kubernetes resource objects. You can use `kubectl explain` to explore the structure of resource objects and determine the JSONPath expressions for fields you want to display.

4. **Examples of Filtering Output with Custom Columns:**

   #### Example 1: Display Pod Names and Node Name
   To display only the names of pods along with their associated node names:
   ```bash
   kubectl get pods -o custom-columns=POD_NAME:.metadata.name,NODE_NAME:.spec.nodeName
   ```
   Output:
   ```
   POD_NAME        NODE_NAME
   pod1            node1
   pod2            node2
   ```

   #### Example 2: Show Service Names and Cluster IP
   Display names of services along with their cluster IPs:
   ```bash
   kubectl get services -o custom-columns=SERVICE_NAME:.metadata.name,CLUSTER_IP:.spec.clusterIP
   ```
   Output:
   ```
   SERVICE_NAME    CLUSTER_IP
   service1        10.1.2.3
   service2        10.1.2.4
   ```

   #### Example 3: List Deployments with Namespace and Labels
   Show deployments with their namespaces and labels:
   ```bash
   kubectl get deployments -o custom-columns=DEPLOYMENT:.metadata.name,NAMESPACE:.metadata.namespace,LABELS:.metadata.labels.app
   ```
   Output:
   ```
   DEPLOYMENT      NAMESPACE   LABELS
   deployment1     default     app1
   deployment2     default     app2
   ```

    #### Example 4: Here is the command that gives custom output showing the pod name and CPU & memory requests.

    ```bash
    kubectl get pods -o custom-columns='POD_NAME:.metadata.name,CPU_REQUEST:.spec.containers[*].resources.requests.cpu,MEMORY_REQUEST:.spec.containers[*].resources.requests.memory'
    ```

    Explanation of the command:
    - `kubectl get pods`: This fetches information about all pods in the cluster.
    - `-o custom-columns='POD_NAME:.metadata.name,CPU_REQUEST:.spec.containers[*].resources.requests.cpu,MEMORY_REQUEST:.spec.containers[*].resources.requests.memory'`: Specifies custom columns to display.
    - `POD_NAME:.metadata.name`: Defines a column header `POD_NAME` that displays the name of each pod (`metadata.name`).
    - `CPU_REQUEST:.spec.containers[*].resources.requests.cpu`: Defines a column header `CPU_REQUEST` that displays the CPU request (`resources.requests.cpu`) for each container (`[*]` means all containers) within the pod.
    - `MEMORY_REQUEST:.spec.containers[*].resources.requests.memory`: Defines a column header `MEMORY_REQUEST` that displays the memory request (`resources.requests.memory`) for each container (`[*]` means all containers) within the pod.

    Output:
    ```
    POD_NAME        CPU_REQUEST   MEMORY_REQUEST
    pod1            1             512Mi
    pod2            0.5           256Mi
    ```

### Notes:
- Ensure you have permissions to access pod information in the specified namespace or across the cluster, depending on your `kubectl` configuration.
- The JSONPath expressions (`spec.containers[*].resources.requests.cpu` and `spec.containers[*].resources.requests.memory`) extract CPU and memory request values from each container within the pods.
- Adjust the command as needed for different resource types or to include additional fields based on your requirements.

### Tips:
- **JSONPath Basics:** Use `kubectl explain` to understand the structure of Kubernetes resources and discover available fields for JSONPath expressions.
- **Formatting:** Pay attention to spacing and commas in your custom-columns specification to ensure correct output formatting.
- **Field Names:** You can name columns (`<header-name>`) whatever you find descriptive for your needs.

### Summary:
Using `custom-columns` with `kubectl get` provides flexibility in filtering and customizing the output of Kubernetes resource information. It's handy for scripting and automation tasks where specific fields from Kubernetes objects are required without unnecessary details. Tailor the command according to your specific use case to efficiently manage and monitor your Kubernetes resources.