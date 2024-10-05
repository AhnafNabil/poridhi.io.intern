# Bootstrapping the etcd Cluster

Kubernetes components are inherently stateless, and they rely on `etcd` to store all cluster-related data, such as configuration, cluster state, and service discovery information. In this guide, we will set up a two-node `etcd` cluster to provide high availability and secure communication between Kubernetes components. This setup is designed to ensure that the cluster state is consistently stored and available even if one `etcd` node goes down.

![](./images/etcd.drawio.svg)

## Prerequisites

The commands in this lab must be run on each controller instances: 
- `controller-0`
- `controller-1`. 

Login to each controller instance using the `ssh` command. Example:

```sh
ssh controller-0
ssh controller-1
```

We can also set the hostname of each controller:

**1. Controller-0**

```sh
sudo hostnamectl set-hostname controller-0
```

**2. Controller-1**
```sh
sudo hostnamectl set-hostname controller-1
```

>Tip: After setting the hostname, exit the SSH session and re-login to reload the changes.

## Bootstrapping an etcd Cluster Member

### Download and Install the etcd Binaries

`etcd` is a distributed key-value store that stores cluster state and configuration information. Download the official etcd release binaries from the [etcd](https://github.com/etcd-io/etcd) GitHub project:

```sh
wget -q --show-progress --https-only --timestamping \
  "https://github.com/etcd-io/etcd/releases/download/v3.5.16/etcd-v3.5.16-linux-amd64.tar.gz"
```

Extract and install the `etcd` server and the `etcdctl` command line utility:

```sh
tar -xvf etcd-v3.5.16-linux-amd64.tar.gz
sudo mv etcd-v3.5.16-linux-amd64/etcd* /usr/local/bin/
```
This command will install both the etcd server and the etcdctl command-line utility, which is used for managing and interacting with the etcd cluster.

### Configure the etcd Server

Create the necessary directories for etcd configuration and data storage:

```sh
sudo mkdir -p /etc/etcd /var/lib/etcd
sudo chmod 700 /var/lib/etcd
sudo cp ca.pem kubernetes-key.pem kubernetes.pem /etc/etcd/
```

The TLS certificates `(ca.pem, kubernetes-key.pem, kubernetes.pem)` are used to secure communication between etcd nodes and clients. The certificates ensure that only trusted nodes and clients can communicate with the etcd cluster.


### Set Environment Variables for etcd Configuration

The internal IP address of each controller node will be used for etcd node communication and client access. Set the environment variables for each controller node as shown below:


## For Controller-0

### Set `INTERNAL_IP`

If you already know the internal IP of the instance, you can set it like this:

```bash
INTERNAL_IP="10.0.1.10"
export INTERNAL_IP
```

### Set `ETCD_NAME`

If you know the etcd name (which typically matches the hostname or a unique identifier for the instance in the etcd cluster), you can set it like this:

```bash
ETCD_NAME="controller-0"
export ETCD_NAME
```

### Verify the Values
You can verify that the environment variables are set correctly by running:

```bash
echo $INTERNAL_IP
echo $ETCD_NAME
```

## For Controller-1

### Set `INTERNAL_IP`
If you already know the internal IP of the instance, you can set it like this:

```bash
INTERNAL_IP="10.0.1.11"
export INTERNAL_IP
```

### Set `ETCD_NAME`

If you know the etcd name (which typically matches the hostname or a unique identifier for the instance in the etcd cluster), you can set it like this:

```bash
ETCD_NAME="controller-1"
export ETCD_NAME
```

### Verify the Values
You can verify that the environment variables are set correctly by running:
```bash
echo $INTERNAL_IP
echo $ETCD_NAME
```

## Create the `etcd.service` systemd unit file:

The `etcd` service is managed by `systemd` to ensure it starts at boot and restarts automatically if it fails. Create the `etcd.service` file with the following content:

```sh
cat <<EOF | sudo tee /etc/systemd/system/etcd.service
[Unit]
Description=etcd
Documentation=https://github.com/coreos

[Service]
Type=notify
ExecStart=/usr/local/bin/etcd \\
  --name ${ETCD_NAME} \\
  --cert-file=/etc/etcd/kubernetes.pem \\
  --key-file=/etc/etcd/kubernetes-key.pem \\
  --peer-cert-file=/etc/etcd/kubernetes.pem \\
  --peer-key-file=/etc/etcd/kubernetes-key.pem \\
  --trusted-ca-file=/etc/etcd/ca.pem \\
  --peer-trusted-ca-file=/etc/etcd/ca.pem \\
  --peer-client-cert-auth \\
  --client-cert-auth \\
  --initial-advertise-peer-urls https://${INTERNAL_IP}:2380 \\
  --listen-peer-urls https://${INTERNAL_IP}:2380 \\
  --listen-client-urls https://${INTERNAL_IP}:2379,https://127.0.0.1:2379 \\
  --advertise-client-urls https://${INTERNAL_IP}:2379 \\
  --initial-cluster-token etcd-cluster-0 \\
  --initial-cluster controller-0=https://10.0.1.10:2380,controller-1=https://10.0.1.11:2380 \\
  --initial-cluster-state new \\
  --data-dir=/var/lib/etcd
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
```

### Explanation of the Configuration
- **Peer URLs:** The `--initial-advertise-peer-urls` and `--listen-peer-urls` flags specify the communication address and port (2380) for etcd nodes to talk to each other.
- **Client URLs:** The `--listen-client-urls` and --advertise-client-urls flags specify the address and port (2379) for clients (like API servers) to communicate with etcd.
- **TLS Certificates:** The `--cert-file`, `--key-file`, and `--trusted-ca-file` options secure communication between etcd nodes and between clients and etcd.
- **Cluster Configuration:** The `--initial-cluster` flag defines the initial etcd cluster members and their peer URLs.

### Start and Enable the etcd Service

```sh
sudo systemctl daemon-reload
sudo systemctl enable etcd
sudo systemctl start etcd
```

> Important: Remember to run the above commands on each controller node: `controller-0`, `controller-1`.

## Verify the etcd Cluster

Once the etcd service is running on both nodes, verify the cluster status by listing the cluster members:

```sh
sudo ETCDCTL_API=3 etcdctl member list \
  --endpoints=https://127.0.0.1:2379 \
  --cacert=/etc/etcd/ca.pem \
  --cert=/etc/etcd/kubernetes.pem \
  --key=/etc/etcd/kubernetes-key.pem
```

> output

```sh
f9b0e395cb8278dc, started, controller-0, https://10.0.1.10:2380, https://10.0.1.10:2379, false
eecdfcb7e79fc5dd, started, controller-1, https://10.0.1.11:2380, https://10.0.1.11:2379, false
```

This output confirms that both `controller-0` and `controller-1` are part of the etcd cluster and are functioning correctly.

So, You have successfully bootstrapped a `two-node` etcd cluster that provides a foundation for storing Kubernetes cluster state and configuration data.