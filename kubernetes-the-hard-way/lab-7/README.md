# Configuring kubectl for Remote Access

In this lab you will generate a kubeconfig file for the `kubectl` command line utility based on the `admin` user credentials.

> Run the commands in this lab from the same directory used to generate the admin client certificates.(poridhi vs code)

## The Admin Kubernetes Configuration File

Each kubeconfig requires a Kubernetes API Server to connect to. To support high availability the IP address assigned to the external load balancer fronting the Kubernetes API Servers will be used.

Generate a kubeconfig file suitable for authenticating as the `admin` user:

```
KUBERNETES_PUBLIC_ADDRESS=$(aws elbv2 describe-load-balancers \
--load-balancer-arns ${LOAD_BALANCER_ARN} \
--output text --query 'LoadBalancers[].DNSName')

kubectl config set-cluster kubernetes-the-hard-way \
  --certificate-authority=ca.pem \
  --embed-certs=true \
  --server=https://${KUBERNETES_PUBLIC_ADDRESS}:443

kubectl config set-credentials admin \
  --client-certificate=admin.pem \
  --client-key=admin-key.pem

kubectl config set-context kubernetes-the-hard-way \
  --cluster=kubernetes-the-hard-way \
  --user=admin

kubectl config use-context kubernetes-the-hard-way
```

## Verification

Check the version of the remote Kubernetes cluster:

```
kubectl version
```

> output

```
Client Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.0", GitCommit:"cb303e613a121a29364f75cc67d3d580833a7479", GitTreeState:"clean", BuildDate:"2021-04-08T16:31:21Z", GoVersion:"go1.16.1", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"21", GitVersion:"v1.21.0", GitCommit:"cb303e613a121a29364f75cc67d3d580833a7479", GitTreeState:"clean", BuildDate:"2021-04-08T16:25:06Z", GoVersion:"go1.16.1", Compiler:"gc", Platform:"linux/amd64"}
```

List the nodes in the remote Kubernetes cluster:

```
kubectl get nodes
```

> output

```
NAME           STATUS   ROLES    AGE     VERSION
ip-10-0-1-20   Ready    <none>   3m35s   v1.21.0
ip-10-0-1-21   Ready    <none>   3m35s   v1.21.0
ip-10-0-1-22   Ready    <none>   3m35s   v1.21.0
```

Next: [Provisioning Pod Network Routes](11-pod-network-routes.md)


# Provisioning Pod Network Routes

Pods scheduled to a node receive an IP address from the node's Pod CIDR range. At this point pods can not communicate with other pods running on different nodes due to missing network [routes](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Route_Tables.html).

In this lab you will create a route for each worker node that maps the node's Pod CIDR range to the node's internal IP address.

> There are [other ways](https://kubernetes.io/docs/concepts/cluster-administration/networking/#how-to-achieve-this) to implement the Kubernetes networking model.

## The Routing Table and routes

In this section you will gather the information required to create routes in the `kubernetes-the-hard-way` VPC network and use that to create route table entries. 

In production workloads this functionality will be provided by CNI plugins like flannel, calico, amazon-vpc-cni-k8s. Doing this by hand makes it easier to understand what those plugins do behind the scenes.

Print the internal IP address and Pod CIDR range for each worker instance and create route table entries:

```sh
ROUTE_TABLE_ID=$(aws ec2 describe-route-tables \
  --filters "Name=tag:Name,Values=kubernetes" \
  --output text --query 'RouteTables[0].RouteTableId')
```

```sh
# Define the pod_cidr values for each worker instance
declare -A pod_cidr_map
pod_cidr_map["worker-0"]="10.200.0.0/24"
pod_cidr_map["worker-1"]="10.200.1.0/24"

# Iterate over each worker instance
for instance in worker-0 worker-1; do
  instance_id_ip="$(aws ec2 describe-instances \
    --filters "Name=tag:Name,Values=${instance}" \
    --output text --query 'Reservations[].Instances[].[InstanceId,PrivateIpAddress]')"

  instance_id="$(echo "${instance_id_ip}" | cut -f1)"
  instance_ip="$(echo "${instance_id_ip}" | cut -f2)"

  # Fetch the pod_cidr from the predefined map
  pod_cidr="${pod_cidr_map[${instance}]}"

  # Check if pod_cidr is empty
  if [[ -z "${pod_cidr}" ]]; then
    echo "Error: pod_cidr is empty for instance ${instance_id} (${instance_ip})"
    continue
  fi

  echo "${instance_ip} ${pod_cidr}"

  # Create route only if pod_cidr is valid
  aws ec2 create-route \
    --route-table-id "${ROUTE_TABLE_ID}" \
    --destination-cidr-block "${pod_cidr}" \
    --instance-id "${instance_id}"
done
```

ssh -i /root/code/k8s-infra-aws/kubernetes.id_rsa ubuntu@${external_ip} \
 "sudo ETCDCTL_API=3 etcdctl get \
  --endpoints=https://127.0.0.1:2379 \
  --cacert=/etc/etcd/ca.pem \
  --cert=/etc/etcd/kubernetes.pem \
  --key=/etc/etcd/kubernetes-key.pem\
  /registry/secrets/default/kubernetes-the-hard-way | hexdump -C"


<!-- ```sh
for instance in worker-0 worker-1; do
  instance_id_ip="$(aws ec2 describe-instances \
    --filters "Name=tag:Name,Values=${instance}" \
    --output text --query 'Reservations[].Instances[].[InstanceId,PrivateIpAddress]')"
  
  instance_id="$(echo "${instance_id_ip}" | cut -f1)"
  instance_ip="$(echo "${instance_id_ip}" | cut -f2)"

  # Fetch the pod-cidr from user data
  pod_cidr=["10.200.0.0", "10.200.0.1"]

  # Check if pod_cidr is empty
  if [[ -z "${pod_cidr}" ]]; then
    echo "Error: pod_cidr is empty for instance ${instance_id} (${instance_ip})"
    continue
  fi

  echo "${instance_ip} ${pod_cidr}"

  # Create route only if pod_cidr is valid
  aws ec2 create-route \
    --route-table-id "${ROUTE_TABLE_ID}" \
    --destination-cidr-block "${pod_cidr}" \
    --instance-id "${instance_id}"
done
``` -->

> output

```sh
10.0.1.20 10.200.0.0/24
{
    "Return": true
}
10.0.1.21 10.200.1.0/24
{
    "Return": true
}
10.0.1.22 10.200.2.0/24
{
    "Return": true
}
```

## Validate Routes

Validate network routes for each worker instance:

```sh
aws ec2 describe-route-tables \
  --route-table-ids "${ROUTE_TABLE_ID}" \
  --query 'RouteTables[].Routes'
```

> output

```
[
    [
        {
            "DestinationCidrBlock": "10.200.0.0/24",
            "InstanceId": "i-0879fa49c49be1a3e",
            "InstanceOwnerId": "107995894928",
            "NetworkInterfaceId": "eni-0612e82f1247c6282",
            "Origin": "CreateRoute",
            "State": "active"
        },
        {
            "DestinationCidrBlock": "10.200.1.0/24",
            "InstanceId": "i-0db245a70483daa43",
            "InstanceOwnerId": "107995894928",
            "NetworkInterfaceId": "eni-0db39a19f4f3970f8",
            "Origin": "CreateRoute",
            "State": "active"
        },
        {
            "DestinationCidrBlock": "10.200.2.0/24",
            "InstanceId": "i-0b93625175de8ee43",
            "InstanceOwnerId": "107995894928",
            "NetworkInterfaceId": "eni-0cc95f34f747734d3",
            "Origin": "CreateRoute",
            "State": "active"
        },
        {
            "DestinationCidrBlock": "10.0.0.0/16",
            "GatewayId": "local",
            "Origin": "CreateRouteTable",
            "State": "active"
        },
        {
            "DestinationCidrBlock": "0.0.0.0/0",
            "GatewayId": "igw-00d618a99e45fa508",
            "Origin": "CreateRoute",
            "State": "active"
        }
    ]
]
```

Next: [Deploying the DNS Cluster Add-on](12-dns-addon.md)