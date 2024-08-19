# SSH Key-Based Authentication to access Multiple EC2 instances via a Bastion Server

## Overview
SSH (Secure Shell) is a protocol for securely connecting to remote systems over a network. Using public/private key pairs for authentication is a more secure method than password-based authentication.

![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-1.png)

We will work on a scenario like this, suppose we have four EC2 instances: one public-facing **bastion server** and **three private server**. We want to SSH into each private servers.

## Before starting, why Use SSH Keys?

- **Enhanced Security:** Passwords can be vulnerable, especially if they are weak or reused. SSH keys provide stronger authentication.

- **Convenience:** Once set up, SSH keys allow you to log in to a server without needing to type a password.

## Step by Step guide

### Step 1: Create Infrastructure(Optional if you have your servers up and running)

For this setup, we will need a Publicly accessible **Bastion server**, and three private servers. We can create these servers in AWS. We can manually create the servers by login in to the AWS management console or we can create the servers and other necessary resouces using PULUMI or Terraform. 

Here is our overall architecture:

![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image.png)

Steps to create the resources and servers in aws using PULUMI.

1. **Configure AWS CLI in your local machine**

  ```sh
  aws configure
  ```

  ![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-2.png)

2. **Setup PULUMI for your project**

- Login into pulumi with your access token

  ```sh
  pulumi login
  ```

- Create an empty directory ( e.g., `Infra-for-ssh`)

  ```sh
  mkdir Infra-for-ssh
  cd Infra-for-ssh
  ```

- Run the following command to create a new Pulumi project:

  ```sh
  pulumi new aws-javascript
  ```
  - Follow the prompts to set up your project.

  ![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-3.png)

3. **Create a Key Pair for your SERVER's**:

We will use specific keys for SSHing into servers. So, we have to create 4 key pairs as we have 4 servers. This is for more secure connection. You can also create a single key pair as well.

- Bastion-server key generation

    ```sh
    aws ec2 create-key-pair --key-name BastionServer --query 'KeyMaterial' --output text > BastionServer.pem
    chmod 400 BastionServer.pem
    ```
- Private-server1 key generation

    ```sh
    aws ec2 create-key-pair --key-name PrivateServer1 --query 'KeyMaterial' --output text > PrivateServer1.pem
    chmod 400 PrivateServer1.pem
    ```

- Private-server1 key generation

    ```sh
    aws ec2 create-key-pair --key-name PrivateServer2 --query 'KeyMaterial' --output text > PrivateServer2.pem
    chmod 400 PrivateServer2.pem
    ```
- Private-server1 key generation

    ```sh
    aws ec2 create-key-pair --key-name PrivateServer3 --query 'KeyMaterial' --output text > PrivateServer3.pem
    chmod 400 PrivateServer3.pem
    ```

![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-4.png)

**NOTE:** Make sure to set the correct permission for the keys

4. **Write the infrastucture creation code in the `index.js` file**

```js
const pulumi = require("@pulumi/pulumi");
const aws = require("@pulumi/aws");

// Create a VPC
const vpc = new aws.ec2.Vpc("my-vpc", {
    cidrBlock: "10.0.0.0/16",
    tags: {
        Name: "my-vpc"
    }
});

exports.vpcId = vpc.id;

// Create a public subnet
const publicSubnet = new aws.ec2.Subnet("public-subnet", {
    vpcId: vpc.id,
    cidrBlock: "10.0.1.0/24",
    availabilityZone: "ap-southeast-1a",
    mapPublicIpOnLaunch: true,
    tags: {
        Name: "public-subnet"
    }
});

exports.publicSubnetId = publicSubnet.id;

// Create a private subnet
const privateSubnet = new aws.ec2.Subnet("private-subnet", {
    vpcId: vpc.id,
    cidrBlock: "10.0.2.0/24",
    availabilityZone: "ap-southeast-1a",
    tags: {
        Name: "private-subnet"
    }
});

exports.privateSubnetId = privateSubnet.id;

// Create an Internet Gateway
const igw = new aws.ec2.InternetGateway("internet-gateway", {
    vpcId: vpc.id,
    tags: {
        Name: "IGW"
    }
});

exports.igwId = igw.id;

// Create a route table for the public subnet
const publicRouteTable = new aws.ec2.RouteTable("public-route-table", {
    vpcId: vpc.id,
    tags: {
        Name: "rt-public"
    }
});

// Create a route in the route table for the Internet Gateway
const route = new aws.ec2.Route("igw-route", {
    routeTableId: publicRouteTable.id,
    destinationCidrBlock: "0.0.0.0/0",
    gatewayId: igw.id
});

// Associate the route table with the public subnet
const routeTableAssociation = new aws.ec2.RouteTableAssociation("public-route-table-association", {
    subnetId: publicSubnet.id,
    routeTableId: publicRouteTable.id
});

exports.publicRouteTableId = publicRouteTable.id;

// Allocate an Elastic IP for the NAT Gateway
const eip = new aws.ec2.Eip("nat-eip", { vpc: true });

// Create the NAT Gateway
const natGateway = new aws.ec2.NatGateway("nat-gateway", {
    subnetId: publicSubnet.id,
    allocationId: eip.id,
    tags: {
        Name: "NGW"
    }
});

exports.natGatewayId = natGateway.id;

// Create a route table for the private subnet
const privateRouteTable = new aws.ec2.RouteTable("private-route-table", {
    vpcId: vpc.id,
    tags: {
        Name: "rt-private"
    }
});

// Create a route in the route table for the NAT Gateway
const privateRoute = new aws.ec2.Route("nat-route", {
    routeTableId: privateRouteTable.id,
    destinationCidrBlock: "0.0.0.0/0",
    natGatewayId: natGateway.id
});

// Associate the route table with the private subnet
const privateRouteTableAssociation = new aws.ec2.RouteTableAssociation("private-route-table-association", {
    subnetId: privateSubnet.id,
    routeTableId: privateRouteTable.id
});

exports.privateRouteTableId = privateRouteTable.id;

// Create a security group for the public instance (Bastion Server)
const publicSecurityGroup = new aws.ec2.SecurityGroup("public-secgrp", {
    vpcId: vpc.id,
    description: "Enable HTTP and SSH access for public instance",
    ingress: [
        { protocol: "tcp", fromPort: 80, toPort: 80, cidrBlocks: ["0.0.0.0/0"] },
        { protocol: "tcp", fromPort: 22, toPort: 22, cidrBlocks: ["0.0.0.0/0"] }
    ],
    egress: [
        { protocol: "-1", fromPort: 0, toPort: 0, cidrBlocks: ["0.0.0.0/0"] }
    ]
});

// Use the specified Ubuntu 24.04 LTS AMI
const amiId = "ami-060e277c0d4cce553";

// Create an EC2 instance in the public subnet (Bastion Server)
const publicInstance = new aws.ec2.Instance("bastion-instance", {
    instanceType: "t2.micro",
    vpcSecurityGroupIds: [publicSecurityGroup.id],
    ami: amiId,
    subnetId: publicSubnet.id,
    keyName: "BastionServer",
    associatePublicIpAddress: true,
    tags: {
        Name: "Bastion-Server"
    }
});

exports.publicInstanceId = publicInstance.id;
exports.publicInstanceIp = publicInstance.publicIp;

// Create a security group for the private instances allowing SSH only from the bastion server
const privateSecurityGroup = new aws.ec2.SecurityGroup("private-secgrp", {
    vpcId: vpc.id,
    description: "Allow SSH access from Bastion server",
    ingress: [
        { protocol: "tcp", fromPort: 22, toPort: 22, cidrBlocks: [pulumi.interpolate`${publicInstance.privateIp}/32`] }
    ],
    egress: [
        { protocol: "-1", fromPort: 0, toPort: 0, cidrBlocks: ["0.0.0.0/0"] }
    ]
});

// Create EC2 instances in the private subnet
const privateInstance1 = new aws.ec2.Instance("private-instance1", {
    instanceType: "t2.micro",
    vpcSecurityGroupIds: [privateSecurityGroup.id],
    ami: amiId,
    subnetId: privateSubnet.id,
    keyName: "PrivateServer1",
    tags: {
        Name: "Private-server1"
    }
});

exports.privateInstance1Id = privateInstance1.id;
exports.privateInstance1PrivateIp = privateInstance1.privateIp;

const privateInstance2 = new aws.ec2.Instance("private-instance2", {
    instanceType: "t2.micro",
    vpcSecurityGroupIds: [privateSecurityGroup.id],
    ami: amiId,
    subnetId: privateSubnet.id,
    keyName: "PrivateServer2",
    tags: {
        Name: "Private-server2"
    }
});

exports.privateInstance2Id = privateInstance2.id;
exports.privateInstance2PrivateIp = privateInstance2.privateIp;

const privateInstance3 = new aws.ec2.Instance("private-instance3", {
    instanceType: "t2.micro",
    vpcSecurityGroupIds: [privateSecurityGroup.id],
    ami: amiId,
    subnetId: privateSubnet.id,
    keyName: "PrivateServer3",
    tags: {
        Name: "Private-server3"
    }
});

exports.privateInstance3Id = privateInstance3.id;
exports.privateInstance3PrivateIp = privateInstance3.privateIp;
```

6. **Deploy the Pulumi Stack:**

    ```sh
    pulumi up
    ```
  - Review the changes and confirm by typing "yes".

![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-5.png)

7. Check the Pulumi output for successfull creation of the Infrastructure


## Step 2: Traditional way: SSH into the private instances via Bastion server

When using the direct SSH command, you have to typically do something like this

1. **First, SSH into the Bastion Server:**

   ```bash
   ssh -i BastionServer ubuntu@<bastion-public-ip>
   ```

2. **Copy the Key files into Bastion Server:**

    ```sh
    scp -i BastionServer.pem PrivateServer1.pem PrivateServer2.pem PrivateServer3.pem privateubuntu@<bastion-public-ip>:~/.ssh/
    ```

![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-6.png)

This command will securely copy the key files into Bastion server using the bastion server key file. After copying the file make sure to set the correct file permission.

3. **Change the file permission of the key files:**

    ```bash
    chmod 400 PrivateServer1.pem
    chmod 400 PrivateServer2.pem
    chmod 400 PrivateServer3.pem
    ```

4. **Then, SSH from the Bastion Server to a Private Instances:**

   ```bash
   ssh -i ~/.ssh/PrivateServer1.pem ubuntu@<private-instance1-ip>
   ssh -i ~/.ssh/PrivateServer2.pem ubuntu@<private-instance2-ip>
   ssh -i ~/.ssh/PrivateServer3.pem ubuntu@<private-instance3-ip>
   ```
![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-7.png)

### Issues with This Approach:

- **Repetitive Commands:** Every time you want to SSH into a private instance, you have to manually SSH into the bastion server first, then SSH into the private instance. This involves typing long commands repeatedly.
  
- **Private Key Management:** You must specify the private key with the `-i` option every time, which can be **cumbersome** if you manage multiple keys for different instances.
  
- **Multi-Hop SSH:** SSHing into private instances requires a multi-hop process, where you first connect to the bastion server and then hop to the private instances. This is not only tedious but also error-prone.

## Simplifying with an SSH Config File

You can solve these issues by configuring the `~/.ssh/config` file. This file allows you to define shortcuts and advanced SSH options, making the SSH process smoother and more efficient.

### Step 1: Set Up the SSH Config File

Here’s an example `~/.ssh/config` file that simplifies the SSH process for this scenario:

```sh
Host bastion
    HostName 18.143.101.33
    User ubuntu
    IdentityFile /root/code/Infra-for-ssh/BastionServer.pem

Host private-server1
    HostName 10.0.2.91
    User ubuntu
    IdentityFile /root/code/Infra-for-ssh/PrivateServer1.pem
    ProxyJump bastion

Host private-server2
    HostName 10.0.2.225
    User ubuntu
    IdentityFile /root/code/Infra-for-ssh/PrivateServer2.pem
    ProxyJump bastion

Host private-server3
    HostName 10.0.2.202
    User ubuntu
    IdentityFile /root/code/Infra-for-ssh/PrivateServer2.pem
    ProxyJump bastion
```
![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-8.png)

**NOTE:**
- Configure the `HostName` with correct IP or DNS name.
- User is `ubuntu` by default if your have launched an ubuntu instance. Change accordingly to your requirement.
- Make sure to change the location of your key files accordingly.

#### Explanation of the Config File:

- **Host bastion:** This section defines the connection to your bastion server. The `HostName` is the public IP of the bastion server, and `IdentityFile` specifies the private key used for authentication.

- **Host private-instance-X:** These sections define the connection to each private instance. The `ProxyJump bastion` directive tells SSH to first connect to the **bastion** server before connecting to the private instance.

### Step 2: Simplified SSH Access

With the above configuration, you can now SSH into any of your instances with simple commands from your **local machine**:

- **SSH into the Bastion Server:**

  ```bash
  ssh bastion
  ```

  ![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-9.png)

- **SSH into Private Server 1:**

  ```bash
  ssh private-server1
  ```

  ![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-10.png)

- **SSH into Private Server 2:**

  ```bash
  ssh private-server2
  ```

- **SSH into Private Server 3:**

  ```bash
  ssh private-server3
  ```

So, we have successfully done our SSH using ssh config file. 

![alt text](https://github.com/Konami33/poridhi.io.intern/raw/main/SSH-Basic/images/image-11.png)

### Benefits of Using an SSH Config File

- **Streamlined Workflow:** No need to remember and type long SSH commands. Just use the server name you defined in the config file.
  
- **Multi-Hop SSH Made Easy:** The `ProxyJump` option allows for seamless SSH connections through the bastion server, making multi-hop SSH effortless.

- **Centralized Key Management:** You don’t need to specify the private key file with each command; it’s all handled by the config file.

### Security Note

- Ensure that your `~/.ssh/config` file is secure by setting appropriate permissions:

    ```bash
    chmod 600 ~/.ssh/config
    ```
This will prevent unauthorized users from viewing or modifying your SSH configuration.

- If you want to push your project into github, make sure to write the keyfiles in the `.gitignore` file.

### Conclusion

By using the SSH config file, you can significantly simplify and secure your SSH access to multiple instances, especially when dealing with bastion servers and private instances. This approach not only saves time but also reduces the potential for errors in your SSH workflow.