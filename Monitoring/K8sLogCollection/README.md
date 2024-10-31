# Kubernetes log collection using Grafana Loki

This lab covers the setup and configuration for collecting Kubernetes logs using Grafana Loki. Loki, paired with Promtail and Grafana, provides an efficient solution for aggregating and visualizing logs from Kubernetes clusters. This documentation will guide you through deploying Loki, Promtail, and Grafana in your Kubernetes environment, configuring them to capture logs, and using Grafana for log analysis.

![](./images/loki-k8s.drawio.svg)

---

## **Overview**

- **Loki**: A log aggregation system designed for simplicity and scalability. Unlike traditional log systems, Loki doesn’t index full logs but relies on a unique approach using labels, which is both cost-effective and fast.
- **Promtail**: A log collector that works as a daemon on each Kubernetes node. It forwards logs from each pod to Loki.
- **Grafana**: A visualization tool that integrates with Loki to provide dashboards for log analysis.

---

## **Prerequisites**

- A Kubernetes cluster with at least one master and worker node.
- **Helm**: Used to simplify the deployment of Loki and Grafana. Install Helm if not already available:
    ```bash
    curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null
    sudo apt-get install apt-transport-https --yes
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
    sudo apt-get update
    sudo apt-get install helm
    ```

    ![alt text](image.png)


## **Deploying Loki, Promtail, and Grafana**

1. **Add Grafana’s Helm Repository**:
   To get Loki, Promtail, and Grafana, add Grafana’s official Helm chart repository.
   ```bash
   helm repo add grafana https://grafana.github.io/helm-charts
   helm repo update
   ```

2. **Explore the Loki Helm Chart Options**:
   Check available Loki stack configurations:
   ```bash
   helm search repo loki
   ```
   ![alt text](image-2.png)

   The preferred option for deploying Loki, Promtail, and Grafana is `loki-stack`.

   ![alt text](image-1.png)

3. **Customize Configuration**:
   Before deploying, check and modify Helm chart values to suit your environment by generating a configuration file:
   ```bash
   helm show values grafana/loki-stack > loki-values.yaml
   ```
   **Modifications**:
   - Enable Grafana by setting `grafana.enabled` to `true`.
   - Set the `grafana.tag` to `latest` for the latest version.


4. **Install Loki Stack**:
   Install Loki with Promtail and Grafana:
   ```bash
   helm install loki grafana/loki-stack -f loki-values.yaml
   ```
   After installation, you should see confirmation that Loki Stack is deployed.

   ![alt text](image-3.png)

5. **Verify Deployment**:
   Check that the required pods are running:
   ```bash
   kubectl get all
   ```
   ![alt text](image-4.png)

---

## **Configuring Grafana**

### **Access Grafana**:
Since the service type is `ClusterIP`, we can not access it from outside of the kubernetes cluster.
To expose the `loki-grafana` service as a NodePort, follow these steps:

1. **Edit the Service**:
   Change the `loki-grafana` service type from `ClusterIP` to `NodePort` to make it accessible from outside the cluster. You can do this by editing the service configuration.

   ```bash
   kubectl edit svc loki-grafana
   ```

2. **Modify the Type and Set the NodePort**:
   In the editor that opens, locate the `spec.type` field and change it to `NodePort`. Then, under the `ports` section, add a `nodePort` field with your desired port (or let Kubernetes assign a random port in the NodePort range if you leave it out).

   Here’s an example:
   ```yaml
   spec:
     type: NodePort
     ports:
       - port: 80
         targetPort: 3000
         protocol: TCP
         nodePort: 32000  # Choose a port within the 30000–32767 range
   ```

3. **Save and Exit**:
   Save the changes and exit the editor. Kubernetes will update the service.

4. **Verify the NodePort Service**:
   Confirm that the `loki-grafana` service is now using a `NodePort` and check the assigned port:
   ```bash
   kubectl get svc loki-grafana
   ```

   You should see the `loki-grafana` service listed as `NodePort` with the specified port number under `PORT(S)`.

5. **Access Grafana**:
   Now, you can access Grafana using the `NodeIP:NodePort` URL:
   ```sh
   http://<node-ip>:32000
   ```

   Replace `<node-ip>` with the IP of any node in your cluster.

This should make Grafana accessible from outside the cluster via the specified NodePort.
   

### **Grafana Credentials**:
Default Grafana username is `admin`. Retrieve the password from the Kubernetes secrets:

```bash
kubectl get secret --namespace <namespace> loki-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```
![alt text](image-5.png)

### **Add Loki as a Data Source**:
   Grafana automatically configures Loki as a data source. Verify this in Grafana under `Settings > Data Sources`.

   ![alt text](image-6.png)

---

## **Log Generator Application**

Here’s a simple log-generating application that can be deployed in Kubernetes. This app will generate **dummy** logs that can be collected by Promtail and viewed in Grafana Loki. The example below is a Python-based application that produces log messages at regular intervals.

### 1. **Application Code** (`app.py`)

Create a directory for example log-application. Then create a file named `app.py` which which will generate log messages with random data.

```python
# app.py
import time
import logging
import random

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger("SimpleLogApp")

# Function to simulate log messages
def generate_log():
    log_levels = [logging.INFO, logging.WARNING, logging.ERROR]
    messages = [
        "User accessed the home page",
        "Database connection established",
        "Failed to retrieve data",
        "API call received",
        "Service temporarily unavailable",
        "User session expired"
    ]
    while True:
        level = random.choice(log_levels)
        message = random.choice(messages)
        logger.log(level, message)
        time.sleep(2)  # Generate log every 2 seconds

if __name__ == "__main__":
    generate_log()
```

### 2. **Dockerfile**

To run this app in a container, create a `Dockerfile`:

```Dockerfile
# Dockerfile
FROM python:3.8-slim

WORKDIR /app

# Copy the Python script
COPY app.py /app

CMD ["python", "app.py"]
```

### 3. **Build and Push Docker Image**

Build and push the Docker image to a registry (like Docker Hub):

```bash
# Build the Docker image
docker build -t <your-docker-username>/simple-log-app:latest .

# Push the Docker image to Docker Hub
docker push <your-docker-username>/simple-log-app:latest
```

![alt text](image-7.png)

Replace `<your-docker-username>` with your Docker Hub username.

### 4. **Kubernetes Deployment** (`app-deployment.yaml`)

Deploy the application on Kubernetes with the following deployment file:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: simple-log-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: simple-log-app
  template:
    metadata:
      labels:
        app: simple-log-app
    spec:
      containers:
      - name: simple-log-app
        image: <your-docker-username>/simple-log-app:latest
        imagePullPolicy: Always
```

Replace `<your-docker-username>` with your actual Docker Hub username.

### 5. **Deploy the Application**

Use `kubectl` to apply the deployment:

```bash
kubectl apply -f app-deployment.yaml
```

![alt text](image-8.png)


## **Promtail Configurations**:

Promtail collects logs from all pods on each node and forwards them to Loki. Promtail is deployed using a DaemonSet to ensure a Promtail instance on every node.

- **Pipeline Stages**: Promtail pipelines allow custom parsing and label extraction. Modify `promtail.yaml` to create labels for specific log attributes.
- **Example**: To extract HTTP method and status code labels from JSON logs.

1. **Edit Promtail Configuration**:
   Retrieve and modify Promtail configuration from the Helm chart secret:
   ```bash
   kubectl get secret loki-promtail -o jsonpath="{.data.promtail\.yaml}" | base64 --decode > promtail.yaml
   ```

2. **Add Pipeline Stages**:
   In `promtail.yaml`, add the following under `pipeline_stages`:
   ```yaml
   pipeline_stages:
     - match:
         selector: '{app="simple-log-app"}'
         stages:
           - json:
               expressions:
                 method: log.method
                 code: log.code
           - labels:
               method:
               code:
   ```

3. **Apply Changes**:
   Update the Promtail configuration in Kubernetes:

    - Generate the Updated YAML:
    ```bash
    kubectl create secret generic loki-promtail --from-file=promtail.yaml --dry-run=client -o yaml > updated-promtail-secret.yaml
    ```
    - Apply the Update:
    ```bash
    kubectl apply -f updated-promtail-secret.yaml
    ```
    ![alt text](image-9.png)

    After updating, restart the Promtail pods to load the new secret.

4. **Check Logs in Grafana**:
   - Go to **Explore** in Grafana.
   - Run a query on the `simple-log-app` pod or app.

   ![alt text](image-10.png)

   - You will see the generated logs:

   ![alt text](image-11.png)
---

## **Conclusion**

With Loki, Promtail, and Grafana deployed on Kubernetes, you have a robust, scalable log aggregation setup that enables efficient log collection and visualization. Promtail's label customization enhances search and filtering, allowing quick insights into application and system logs.