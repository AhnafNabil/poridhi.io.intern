# Bootstrapping the Kubernetes Worker Nodes

In this lab you will learn how to bootstrap two Kubernetes worker nodes, install necessary components, and configure them to join the Kubernetes cluster. The worker nodes are where your applications and services will run. We will configure the worker nodes to communicate with the control plane, manage containers, and handle networking.

![](./images/worker-1.drawio.svg)


## Overview
The worker nodes in a Kubernetes cluster run critical components that enable them to host containers, manage networking, and communicate with the control plane. Each worker node is provisioned with the following components:

- **runc:** A lightweight container runtime for managing container processes.
- **Container Networking Interface (CNI) Plugins:** Provide networking capabilities such as setting up routes and IP addresses.
- **containerd:** Manages container lifecycle and ensures containers are running as expected.
- **kubelet:** The primary Kubernetes node agent responsible for managing pods and container health.
- **kube-proxy:** Handles networking and communication between services within the cluster.

This guide covers bootstrapping two worker nodes and configuring them to communicate securely with the control plane.

## Prerequisites

The commands in this lab must be run on each worker instance: `worker-0`, `worker-1`. Login to each worker instance using the `ssh` command.

## Step 01: Install the OS Dependencies

Start by installing the necessary OS dependencies on each worker node. This includes tools like `socat` and `conntrack`, which are required for Kubernetes operations.

```sh
sudo apt-get update
sudo apt-get -y install socat conntrack ipset
```

- **socat:** Enables support for the `kubectl port-forward` command.
- **conntrack**: Manages connections to ensure that packet flows are correctly tracked and maintained.

## Step 02: Disable Swap

Kubernetes requires that swap be disabled to ensure proper memory allocation and resource management. Check if swap is enabled:

```sh
sudo swapon --show
```

If output is empthy then swap is not enabled. If swap is enabled run the following command to disable swap immediately:

```sh
sudo swapoff -a
```

> NOTE: To ensure swap remains off after reboot consult your Linux distro documentation.

## Step 03: Download and Install Worker Binaries

Download the required binaries for Kubernetes worker nodes:

```sh
wget -q --show-progress --https-only --timestamping \
  https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.21.0/crictl-v1.21.0-linux-amd64.tar.gz \
  https://github.com/opencontainers/runc/releases/download/v1.0.0-rc93/runc.amd64 \
  https://github.com/containernetworking/plugins/releases/download/v0.9.1/cni-plugins-linux-amd64-v0.9.1.tgz \
  https://github.com/containerd/containerd/releases/download/v1.4.4/containerd-1.4.4-linux-amd64.tar.gz \
  https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kubectl \
  https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kube-proxy \
  https://storage.googleapis.com/kubernetes-release/release/v1.21.0/bin/linux/amd64/kubelet
```

Create the directories required for storing these binaries and configurations:

```sh
sudo mkdir -p \
  /etc/cni/net.d \
  /opt/cni/bin \
  /var/lib/kubelet \
  /var/lib/kube-proxy \
  /var/lib/kubernetes \
  /var/run/kubernetes
```

Extract and install the downloaded binaries:

```sh
mkdir containerd
tar -xvf crictl-v1.21.0-linux-amd64.tar.gz
tar -xvf containerd-1.4.4-linux-amd64.tar.gz -C containerd
sudo tar -xvf cni-plugins-linux-amd64-v0.9.1.tgz -C /opt/cni/bin/
sudo mv runc.amd64 runc
chmod +x crictl kubectl kube-proxy kubelet runc 
sudo mv crictl kubectl kube-proxy kubelet runc /usr/local/bin/
sudo mv containerd/bin/* /bin/
```

## Step 4: Configure Container Networking Interface (CNI)

The CNI plugins are used to set up container networking and handle IP address management. Retrieve the Pod CIDR range for the current compute instance:

```sh
POD_CIDR="10.200.0.0/16"
export POD_CIDR
echo $POD_CIDR
```

Create the bridge network configuration file for CNI:

```sh
cat <<EOF | sudo tee /etc/cni/net.d/10-bridge.conf
{
    "cniVersion": "0.4.0",
    "name": "bridge",
    "type": "bridge",
    "bridge": "cnio0",
    "isGateway": true,
    "ipMasq": true,
    "ipam": {
        "type": "host-local",
        "ranges": [
          [{"subnet": "${POD_CIDR}"}]
        ],
        "routes": [{"dst": "0.0.0.0/0"}]
    }
}
EOF
```

Create the loopback network configuration file:

```sh
cat <<EOF | sudo tee /etc/cni/net.d/99-loopback.conf
{
    "cniVersion": "0.4.0",
    "name": "lo",
    "type": "loopback"
}
EOF
```

These configurations ensure that the pods can communicate within the node and across nodes using the bridge and loopback interfaces.

## Step 5: Configure containerd

`containerd` is a high-performance container runtime that interacts with `runc`. Create the `containerd` configuration file:

```sh
sudo mkdir -p /etc/containerd/
```

Then, add the configuration details:

```sh
cat > config.toml << EOF
[plugins]
  [plugins.cri.containerd]
    snapshotter = "overlayfs"
    [plugins.cri.containerd.default_runtime]
      runtime_type = "io.containerd.runc.v2"
      runtime_engine = "/usr/local/bin/runc"
      runtime_root = ""
EOF
```

Create the `containerd.service` systemd unit file:

```sh
cat <<EOF | sudo tee /etc/systemd/system/containerd.service
[Unit]
Description=containerd container runtime
Documentation=https://containerd.io
After=network.target

[Service]
ExecStartPre=/sbin/modprobe overlay
ExecStart=/usr/bin/containerd
Restart=always
RestartSec=5
Delegate=yes
KillMode=process
OOMScoreAdjust=-999
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
CPUAccounting=true
MemoryAccounting=true

[Install]
WantedBy=multi-user.target
EOF
```

## Step 06: Configure the Kubelet

Move the necessary files for the kubelet, including the node-specific certificates and kubeconfig:

### For worker-0:

```sh
WORKER_NAME="worker-0"
echo "${WORKER_NAME}"

sudo mv ${WORKER_NAME}-key.pem ${WORKER_NAME}.pem /var/lib/kubelet/
sudo mv ${WORKER_NAME}.kubeconfig /var/lib/kubelet/kubeconfig
sudo mv ca.pem /var/lib/kubernetes/
```

### For worker-1:

```sh
WORKER_NAME="worker-1"
echo "${WORKER_NAME}"

sudo mv ${WORKER_NAME}-key.pem ${WORKER_NAME}.pem /var/lib/kubelet/
sudo mv ${WORKER_NAME}.kubeconfig /var/lib/kubelet/kubeconfig
sudo mv ca.pem /var/lib/kubernetes/
```

Create the `kubelet-config.yaml` configuration file:

```sh
cat <<EOF | sudo tee /var/lib/kubelet/kubelet-config.yaml
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
authentication:
  anonymous:
    enabled: false
  webhook:
    enabled: true
  x509:
    clientCAFile: "/var/lib/kubernetes/ca.pem"
authorization:
  mode: Webhook
clusterDomain: "cluster.local"
clusterDNS:
  - "10.32.0.10"
podCIDR: "10.200.1.0/24"  # Replace with the actual Pod CIDR for this worker
resolvConf: "/run/systemd/resolve/resolv.conf"
runtimeRequestTimeout: "15m"
tlsCertFile: "/var/lib/kubelet/${WORKER_NAME}.pem"
tlsPrivateKeyFile: "/var/lib/kubelet/${WORKER_NAME}-key.pem"
cgroupDriver: systemd
startupGracePeriod: "30s"
EOF
```

> The `resolvConf` configuration is used to avoid loops when using CoreDNS for service discovery on systems running `systemd-resolved`. 

Create the `kubelet.service` systemd unit file:

### worker-0

```sh
cat <<EOF | sudo tee /etc/systemd/system/kubelet.service
[Unit]
Description=Kubernetes Kubelet
Documentation=https://github.com/kubernetes/kubernetes
After=containerd.service
Requires=containerd.service

[Service]
ExecStart=/usr/local/bin/kubelet \\
  --config=/var/lib/kubelet/kubelet-config.yaml \\
  --container-runtime=remote \\
  --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock \\
  --image-pull-progress-deadline=2m \\
  --kubeconfig=/var/lib/kubelet/kubeconfig \\
  --network-plugin=cni \\
  --register-node=true \\
  --v=2
Restart=on-failure
RestartSec=5
Delegate=yes
KillMode=process
OOMScoreAdjust=-999
TasksMax=infinity

[Install]
WantedBy=multi-user.target
EOF
```

### worker-1

```sh
cat <<EOF | sudo tee /etc/systemd/system/kubelet.service
[Unit]
Description=Kubernetes Kubelet
Documentation=https://github.com/kubernetes/kubernetes
After=containerd.service
Requires=containerd.service

[Service]
ExecStart=/usr/local/bin/kubelet \\
  --config=/var/lib/kubelet/kubelet-config.yaml \\
  --container-runtime=remote \\
  --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock \\
  --image-pull-progress-deadline=2m \\
  --kubeconfig=/var/lib/kubelet/kubeconfig \\
  --network-plugin=cni \\
  --register-node=true \\
  --v=2
Restart=on-failure
RestartSec=5
Delegate=yes
KillMode=process
OOMScoreAdjust=-999
TasksMax=infinity

[Install]
WantedBy=multi-user.target
EOF
```

## Step 07: Configure the Kubernetes Proxy

Move the kube-proxy configuration file:

```sh
sudo mv kube-proxy.kubeconfig /var/lib/kube-proxy/kubeconfig
```

Create the `kube-proxy-config.yaml` configuration file:

```sh
cat <<EOF | sudo tee /var/lib/kube-proxy/kube-proxy-config.yaml
kind: KubeProxyConfiguration
apiVersion: kubeproxy.config.k8s.io/v1alpha1
clientConnection:
  kubeconfig: "/var/lib/kube-proxy/kubeconfig"
mode: "iptables"
clusterCIDR: "10.200.0.0/16"
iptables:
  syncPeriod: "30s"
EOF
```

Create the `kube-proxy.service` systemd unit file:

```sh
cat <<EOF | sudo tee /etc/systemd/system/kube-proxy.service
[Unit]
Description=Kubernetes Kube Proxy
Documentation=https://github.com/kubernetes/kubernetes
After=network.target

[Service]
ExecStart=/usr/local/bin/kube-proxy \\
  --config=/var/lib/kube-proxy/kube-proxy-config.yaml
Restart=on-failure
RestartSec=5
Delegate=yes
KillMode=process
OOMScoreAdjust=-999
TasksMax=infinity

[Install]
WantedBy=multi-user.target
EOF
```

## Step 08: Start the Worker Services
Enable and start all the required services on each worker node:

```sh
sudo systemctl daemon-reload
sudo systemctl enable containerd kubelet kube-proxy
sudo systemctl start containerd kubelet kube-proxy
```

> Remember to run the above commands on each worker node: `worker-0`, `worker-1`.

## Verification

> The compute instances created in this tutorial will not have permission to complete this section. Run the following commands from the same machine used to create the compute instances.

List the registered Kubernetes nodes:

```sh
external_ip=$(aws ec2 describe-instances --filters \
    "Name=tag:Name,Values=controller-0" \
    "Name=instance-state-name,Values=running" \
    --output text --query 'Reservations[].Instances[].PublicIpAddress')

ssh -i ~/.ssh/kubernetes.id_rsa ubuntu@${external_ip} kubectl get nodes --kubeconfig admin.kubeconfig
```

> output

```
NAME             STATUS   ROLES    AGE   VERSION
ip-10-0-1-20   Ready    <none>   51s   v1.21.0
ip-10-0-1-21   Ready    <none>   51s   v1.21.0
```

This output indicates that all worker nodes have successfully joined the Kubernetes cluster and are in the "Ready" state.


**Congratulations!** You have successfully bootstrapped the Kubernetes worker nodes and configured them to join your cluster. In the next steps, you can configure kubectl for remote access and start deploying applications on the worker nodes.