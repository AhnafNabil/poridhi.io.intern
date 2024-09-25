# Bootstrapping the etcd Cluster

Kubernetes components are stateless and store cluster state in [etcd](https://github.com/etcd-io/etcd). In this lab you will bootstrap a three node etcd cluster and configure it for high availability and secure remote access.

## Prerequisites

The commands in this lab must be run on each controller instance: `controller-0`, `controller-1`. Login to each controller instance using the `ssh` command. Example:

```sh
ssh controller-0
ssh controller-1
```

```sh
sudo hostnamectl set-hostname controller-0
sudo hostnamectl set-hostname controller-1
```

Now ssh into each one of the IP addresses received in last step.

## Bootstrapping an etcd Cluster Member

### Download and Install the etcd Binaries

Download the official etcd release binaries from the [etcd](https://github.com/etcd-io/etcd) GitHub project:

```sh
wget -q --show-progress --https-only --timestamping \
  "https://github.com/etcd-io/etcd/releases/download/v3.5.16/etcd-v3.5.16-linux-amd64.tar.gz"
```

Extract and install the `etcd` server and the `etcdctl` command line utility:

```sh
tar -xvf etcd-v3.5.16-linux-amd64.tar.gz
sudo mv etcd-v3.5.16-linux-amd64/etcd* /usr/local/bin/
```

### Configure the etcd Server

```sh
sudo mkdir -p /etc/etcd /var/lib/etcd
sudo chmod 700 /var/lib/etcd
sudo cp ca.pem kubernetes-key.pem kubernetes.pem /etc/etcd/
```

The instance internal IP address will be used to serve client requests and communicate with etcd cluster peers. Retrieve the internal IP address for the current compute instance:


## Controller-0

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

## Controller-1

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

### Start the etcd Server

```sh
sudo systemctl daemon-reload
sudo systemctl enable etcd
sudo systemctl start etcd
```

> Remember to run the above commands on each controller node: `controller-0`, `controller-1`.

## Verification

List the etcd cluster members:

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

Next: [Bootstrapping the Kubernetes Control Plane](08-bootstrapping-kubernetes-controllers.md)