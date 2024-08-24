# Ray Labs: Infra Deployment with Pulumi + Manual deployment of ray clusters + XGBoost Model Tranining

## Description/Introduction:

Welcome to this  lab where we will deep dive into infrastructure deployment using Pulumi, an open-source infrastructure as code project. We will not just stop there but also manually deploy Ray clusters. Ray is a high-performance distributed execution framework targeted at large-scale machine learning and reinforcement learning applications.

For our practical use case, we will be running a simple linear regression model. Linear regression is a fundamental algorithm in machine learning, making it a perfect fit to demonstrate the capabilities of our deployed infrastructure and Ray clusters. This lab aims to provide a comprehensive understanding of how to deploy, manage, and utilize cloud resources for machine learning tasks effectively.

By the end of this lab, you will have hands-on experience with Pulumi, Ray clusters, and running machine learning models on the cloud, equipping you with valuable skills for your cloud and machine learning journey.

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image.png)

## MLOps Objectives of this Lab:

- Streamline ML workflow
- Automate processes
- Ensure reproducibility
- Monitor and improve models
- Collaborate and version control
- Address scalability and security

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-1.png)

## Setting up the Local Environment:

It would be better if you have a linux/ubuntu OS running in your local machine. If you are on Windows, please setup WSL2 to get the most out of this lab. 

### **1. Install the AWS CLI**

1. set it with aws configure
2. show the screenshots from the terminal

To install AWS CLI using pip and configure it, follow these steps:

### Step 1: Install pip (if not already installed)

1. **Open a terminal session**:
2. **Install pip**:
    - Run the command:
        
        ```bash
        sudo apt install python3-pip
        ```
        
    - This command installs pip for Python 3 on Ubuntu-based systems.

### Step 2: Install AWS CLI using pip

1. **Install AWS CLI**:
    - Run the command:
        
        ```bash
        python3 -m pip install awscli
        ```
        
    - This command installs the AWS CLI and its dependencies.

### Step 3: Verify the Installation

1. **Check the AWS CLI version**:
    - Run the command:
        
        ```bash
        aws --version
        ```
        
    - This command outputs the version of the AWS CLI installed, confirming that the installation was successful.

### Step 4: Configure AWS CLI

1. **Configure AWS CLI**:
    - Run the command:
        
        ```bash
        aws configure
        ```
        
    - This command prompts you for your AWS Access Key ID, Secret Access Key, region, and output format.

### Step 5: Verify AWS CLI Configuration

1. **Verify AWS CLI configuration**:
    - Run the command:
        
        ```bash
        aws ec2 describe-instances
        ```
        
    - This command lists the instances running in the specified region. You have not created any instances by far, you will see an empty list. 

## **2. Create a Key-pair from AWS Console**

1. Log into the AWS console
2. Search for “Key pairs” 
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-2.png)
    
3. In the Key pair (login) section, click Create new key pair.
4. In the Create key pair window, for the Key pair name, type `raycluster-aws`
5. Leave the Key pair type set to RSA.
6. For the private key file format, choose either .pem or .ppk, depending on if you will use PuTTY or OpenSSH to connect to the EC2 instance later on. Note: Use .pem for a mac device; use .ppk for a Windows device.
7. Click Create key pair.

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-3.png)

h. Your `raycluster-aws.pem` or .ppk key pair file should have downloaded. You will use this file later in the lab.

## **3. Clone the GitHub Repo**

Clone the following repo in your local machine where the script and configuration for setting up the AWS components are included,

https://github.com/tahhnik/raycluster-awsdeployment-pypulumi

[GitHub - tahhnik/xgboost-ray-train-ml-workflow-notebooks](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/https://github.com/tahhnik/xgboost-ray-train-ml-workflow-notebooks/tree/main)

# What is Ray?

Ray is an open-source unified framework for distributed computing, and scaling AI and Python applications like machine learning. It provides the compute layer for parallel processing so that you don’t need to be a distributed systems expert. Ray minimizes the complexity of running your distributed individual and end-to-end machine learning workflows with these components:

- Scalable libraries for common machine learning tasks such as data preprocessing, distributed training, hyperparameter tuning, reinforcement learning, and model serving.
- Pythonic distributed computing primitives for parallelizing and scaling Python applications.
- Integrations and utilities for integrating and deploying a Ray cluster with existing tools and infrastructure such as Kubernetes, AWS, GCP, and Azure.

Ray works by **creating a cluster of nodes** and **scheduling tasks** across them. It uses ***dynamic task*** ***scheduling***, ***a shared memory object store*** for efficient data sharing, supports the actor model for stateful computations, and has native support for Python and Java. It's a powerful tool for distributed computing, providing efficient task scheduling, data sharing, and support for stateful computations.

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-4.png)

In this lab, we’ll create a ray cluster with 3 nodes, 

1. One Head Node,
2. Two Worker Nodes, 

Supported by the compute and networking infrastructure of AWS

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-5.png)

## Infrastructure Needed to set up Ray Cluster and ML Workflow on AWS:

Having the right network configurations and setup in AWS for Ray clusters is crucial for efficient distributed computing.  The compute setup, including the choice of EC2 instances and their configurations, directly impacts the performance and cost-effectiveness of the Ray cluster. 

Moreover, we need robust and dynamic storage support to ingest dataset, store the transformed and computed feature data to a feature store, store the model training and output data.

## Compute Instances Configurations:

1. **EC2 Instances**: AWS EC2 instances are used to host the Ray nodes. The instance type (e.g., t2.xlarge) determines the hardware of the host computer.
2. **AMI (Amazon Machine Image)**: Provides the information required to launch an instance, which is a virtual server in the cloud.
3. **User Data Scripts**: These are used to bootstrap instances, installing necessary software and starting services. In the case of Ray, this might include installing Python, setting up a virtual environment, and installing and starting Ray.
4. **EBS (Elastic Block Store)**: Provides persistent block storage volumes for use with EC2 instances. It's used to store data that should persist beyond the life of the instance.
5. **Key Pair**: Used to securely connect to instances which you have already created and obtained from AWS console. It's crucial for managing the Ray cluster.

## Needed Storage Capabilities:

1. **S3 Buckets:** S3 is an object storage of AWS where data is stored in objects, which have three main components: the object’s content or data, a unique identifier for the object, and descriptive metadata including the object’s name, URL, and size

## Needed Networking Configurations:

1. **VPC (Virtual Private Cloud)**: An isolated network environment in AWS. It's the foundation of the AWS Cloud network.
2. **Internet Gateway**: Connects the VPC to the internet, enabling Ray nodes to communicate with the outside world.
3. **Subnet**: A segment of the VPC's IP address range where you can place groups of isolated resources. Ray nodes are placed in subnets.
4. **Route Table**: Directs traffic in the VPC. It's important to route traffic correctly between nodes and to/from the internet.
5. **Route Table Association**: Associates a subnet with a route table, ensuring that the Ray nodes in the subnet follow the correct routing rules.
6. **Security Group**: Acts as a virtual firewall for controlling inbound and outbound traffic. It's crucial to allow specific traffic (like SSH) for cluster management and inter-node communication.

**Creating these components manually on AWS console are very time-consuming and gruesome tasks.** And, in the context of ML engineering, it’s very necessary to build up your infrastructure fast and focus on the ML tasks. 

ML Engineers/Data Scientists/ML Researchers or MLOps Engineers should focus more on the ML workflows rather than allocating most of their time on setting up infrastructures and system configurations. And, to train, build and serve scalable AI/ML models/systems, practitioners often need support of distributed computing or clusters.

The tools we are discussing in this lab, are meant to meet those needs.

Firsly, *we are going to deploy the network, compute and storage infrastructure needed to support Ray cluster and ML Workflow through a **IaaC SDK named Pulumi.***

Secondly, *we are going to **install and launch Ray cluster** on the head node, and connecting other two EC2 instances as worker nodes.*

Thirdly, ***we’ll train a minimal machine learning model with XGBoost** and other algorithms addressing the **fundamental operations of MLOps** with the help of **Ray AIR Libraries and it’s integrations***

Finally, *we’ll deploy the model and **monitor the inferencing with Ray monitor** and other Ray integrations.*

**Now, at first, let’s build the network, compute and storage configurations and deploy the infrastructure on AWS with Pulumi;**

## Setting up the Infrastructure for Model Training:

### What is Pulumi?

Pulumi is an open-source infrastructure-as-code software that allows users to manage cloud infrastructure resources using programming languages such as Go, JavaScript, TypeScript, Python, Java, C#, and YAML. It supports deployment to various cloud providers like AWS, Azure, Google Cloud, and Kubernetes

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-6.png)

### **Step 1: Install Pulumi**

1. **Open a terminal on your machine**.
2. **Install Pulumi using the installation script**: This command will download and install the latest version of Pulumi. If you want to install a specific version, replace **`v3.120.0`** with your desired version.
    
    `bash$ curl -fsSL https://get.pulumi.com | sh -s -- --version dev`
    

### **Step 2: Verify Installation**

1. **Run the `pulumi` CLI**: This command should display the version of Pulumi installed on your machine.
    
    `bash$ pulumi version
    v3.120.0`
    

### Step 3: Set up Your Pulumi Cloud Account and Connect to Pulumi CLI:

To set up a Pulumi Cloud account and log in to the CLI, follow these steps:

#### Create a Pulumi Cloud Account

1. **Navigate to the Pulumi Cloud Console**:  [***https://app.pulumi.com***](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/https://app.pulumi.com)
2. **Create an account** if you don't already have one.
3. **Create and copy the personal access token**
    
    You can find your account’s personal access tokens on the below section,
    
![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-7.png)

#### Configure Your Pulumi CLI

1. **Run `pulumi login`**:
    
    ```bash
    $ pulumi login
    ```
    
    This command will prompt you to log in to your Pulumi Cloud account.
    
2. **Enter your access token**:
    
    `bash$ Enter your access token from https://app.pulumi.com/account/tokens`
    

#### Verify the Login

1. **Run `pulumi about`**:
    
    ```bash
    $ pulumi about
    ```
    
    This command will display information about your Pulumi CLI, including your backend, user, and organizations.
    

#### Additional Tips

- **Verify your Pulumi account**: Ensure you have a Pulumi account and are logged in. You can do this by running `pulumi login` and following the prompts.
- **Configure your Pulumi stack**: You can configure your Pulumi stack by running `pulumi stack init` and following the prompts.
- **Monitor your deployment**: You can monitor your deployment using the Pulumi Service Console or by running `pulumi stack output` to view the output of your Pulumi program.

### Step 4: Programmatically Set up the AWS Deployment Configurations

#### Create a new Directory

```bash
mkdir pulumi-awsraycluster-deployment
```

change the current working directory to the newly created **`pulumi-awsraycluster-deployment`** directory.

```bash
cd pulumi-awsraycluster-deployment
```

#### Create a New Pulumi Project with ‘aws-python’ template

```bash
pulumi new aws-python
```

If prompted, provide the name for your project, project description and you can leave the stack name as default (by hitting enter). 

Provide the aws:region name “ap-southeast-1” (we’ll be using ap-southeast-1 (Singapore) for the AWS region)

Then it would initiate the project setup and install some dependencies.

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-8.png)

After finishing setup, you can list all the contents in the directory and you can see that the files for newly created Pulumi project is there. 

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-9.png)

#### Setting up the Pulumi Deployment Script

At first, go back to the directory which we have cloned from GitHub at the first place. 

Open the [**`ray-cluster-awsdeployment-config.py`](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/http://ray-cluster-awsdeployment-config.py)** and copy it.

Now, head back to the pulumi project directory and open the `__main__.py` file with your code editor.

Replace the code with previously copied code from **`ray-cluster-awsdeployment-config.py`**

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-10.png)

This Pulumi code creates a VPC, an Internet Gateway, a public subnet, a route table, a security group, a head node, and worker nodes for a Ray cluster. The head node is responsible for starting the Ray cluster, while the worker nodes connect to the cluster using the head node's private IP address. Now, let’s look at the networking and compute configurations  in the Pulumi script **`__main__.py`**,

## Setting up Networking Configurations for The EC2 instances:

### Create a VPC

1. The code starts by creating a VPC using `aws.ec2.Vpc()`. This VPC will have a CIDR block of `10.0.0.0/16`, and it will have DNS support and DNS hostnames enabled.  This will create a VPC with up to 65,536 private IPv4 addresses.
2. The VPC ID is stored in the `vpc.id` variable. DNS support and DNS hostnames for this VPC are enabled by setting `enable_dns_support=True` and `enable_dns_hostnames=True`.
3. **What’s a VPC?:** A private network in AWS. Ray clusters need a VPC to operate in an isolated network environment so that the operations can be secured.

### Create an Internet Gateway

1. Next, the code creates an Internet Gateway using `aws.ec2.InternetGateway()`.
2. The Internet Gateway is attached to the VPC using the `vpc_id` parameter.
3. **Why do we need an IGw in here?:** Connects the VPC to the internet. Necessary for Ray nodes to communicate with the outside world, download dependencies, etc.

### Create a Public Subnet

1. A public subnet is created using `aws.ec2.Subnet()`.
2. The subnet is associated with the VPC using the `vpc_id` parameter.
3. The subnet has a CIDR block of `10.0.1.0/24` and is configured to map public IP addresses on launch.
4. Set `map_public_ip_on_launch=True` to ensure that any instances launched in this subnet will automatically be assigned a public IP address
5. **Why a Subnet is Needed:** A segment of the VPC's IP address range where you can place groups of isolated resources. Ray nodes are placed in subnets.

### Create a Route Table

1. A route table is created using `aws.ec2.RouteTable()`.
2. The route table is associated with the VPC using the `vpc_id` parameter.
3. A route is added to the route table, directing all traffic (`0.0.0.0/0`) to the Internet Gateway id set in step 2.
4. **Why it is important for Ray:** Defines rules to direct traffic in the VPC. For Ray, it's important to route traffic correctly between nodes and to/from the internet

### Associate the Subnet with the Route Table

1. The public subnet is associated with the route table using `aws.ec2.RouteTableAssociation()`.
2. The `subnet_id` and `route_table_id` parameters are used to link the subnet and the route table.
3. **Route Table Association**: Associates a subnet with a route table. This ensures that the Ray nodes in the subnet follow the correct routing rules.

CIDR blocks ("10.0.0.0/16" for VPC and "10.0.1.0/24" for subnet) define the IP address range for the VPC and subnet, respectively. They are crucial for network partitioning and routing.

### Create a Security Group

The security group in the script is configured to allow traffic on specific ports, which are crucial for the operation of a Ray cluster and make it secure to work on distributed AI/ML workloads:

1. A security group is created using `aws.ec2.SecurityGroup()`.
2. The security group is associated with the VPC using the `vpc_id` parameter.
3. **IMPORTANT:** Inbound rules are added to the security group, allowing SSH (port 22), HTTP (port 80), HTTPS (port 443), Redis (ports 6379-6382), and a range of ports (1024-65535).
    - **Port 22**: Used for SSH connections, allowing you to remotely manage and control the nodes in your cluster.
    - **Port 80 and 443**: Typically used for HTTP and HTTPS traffic, which could be necessary for downloading dependencies or data.
    - **Ports 6379 to 6382**: These are used by Redis, which Ray uses as its primary data store.
    - **Ports 1024 to 65535**: These are ephemeral ports that can be used by the Ray processes for inter-node communication. These ports are also necessary to be used by Ray Dashboard, Prometheus, Grafana, MLFlow, Tensorboard and other cluster/infra monitoring, experiment tracking and model monitoring tools.
4. An outbound rule is added, allowing all traffic (`0.0.0.0/0`) to go out.

## Setting up the Configurations for the Head and Worker Nodes (EC2 Instances for Ray Cluster):

### Create the Head Node

1. The code creates the head node (EC2 instance) using `aws.ec2.Instance()`.
2. The head node is configured with the following:
    - Instance type: `t2.medium`
    - AMI: `ami-003c463c8207b4dfa`
    - Security group: the one created in the previous step
    - Subnet: the public subnet created earlier
    - User data: a script that installs Python 3.9, creates a virtual environment, and installs Ray
    - Key pair: `key-pair-poridhi-poc` (replace with your key pair name)
    - EBS block device: a 100 GB GP3 volume

### Create the Worker Nodes

1. The code creates the worker nodes (EC2 instances) using a loop and `aws.ec2.Instance()`.
2. Each worker node is configured with the following:
    - Instance type: `t2.medium`
    - AMI: `ami-003c463c8207b4dfa`
    - Security group: the one created in the previous step
    - Subnet: the public subnet created earlier
    - User data: a script that installs Python 3.9, creates a virtual environment, and installs Ray, just like the head node.
    - Key pair: `key-pair-poridhi-poc` (replace with your key pair name)
    - EBS block device: a 100 GB GP3 volume

## Setting up the Configurations for the S3 Buckets (for staging raw dataset, feature store, model training data store and model output store):

- The code will create four AWS S3 buckets with unique names.
- Each bucket is configured to be private and has versioning enabled.
- The unique names for the buckets are specified directly in the code, add unique names to the places (e.g: your own name or other words)

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-11.png)

### Export Outputs

1. The code exports the private IP address of the head node and the worker nodes as Pulumi outputs.
2. The code also exports the names of the s3 buckets

**You can see the deployment plan by running the following command on terminal,**

(It may take some time to download the necessary plugins)

```bash
pulumi preview
```

As you can see, we have all of the components listed in there ready to be deployed and provisioned,

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-12.png)

### Step 5: Deploy the Infrastructure

As the cloud infrastructure is ready to deploy, use the following command to deploy them on the cloud,

```bash
pulumi up
```

When prompted, select yes and the deployment will start,

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-13.png)

The deployment may take 2-3 minutes to complete. After deployment, it’d show that all the network, compute and storage components are created,

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-14.png)

You can check them out by going back to the AWS console where the networking components, EC2 instances and S3 buckets are created, (search with EC2, S3, VPC to get to those windows)

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-15.png)

You can also check the created the S3 buckets by searching S3 on the search bar and going to the S3 window.

Now, the EC2 instances are ready to be deployed as nodes of the Ray Cluster.

## Setting up The Ray Cluster:

We have set up the AWS components needed to set up the Ray cluster with IaaC but the installation of Ray and its dependencies will be installed manually to contextualize better. Eventually, the head node and the worker nodes of the Ray cluster would be added manually.

### For the Head Node:

1. **Connect the head node from User end**
    
    From the list of instances in the EC2 console, find the instance with the head node’s private ip shown in pulumi export from terminal.  Right-click on the instance and click Connect.On the Connect to instance page, select the SSH client tab.
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-16.png)
    
    1. On the SSH client tab, directly underneath Example:, copy the connection command to your clipboard. (You can optionally paste it into a text editor to keep it handy.)
    
    1. Open your SSH client.
    2. Locate your private key file on your local machine that you downloaded earlier, `raycluster-aws`(or .ppk). If it is in your local Downloads folder, you can run cd ~/Downloads.
    3. Run the command chmod 400 `poridhi-one.pem` , to ensure your key is not publicly viewable.
    4. Connect to the instance by running the connection command you copied to your clipboard earlier.
    5. Type in yes to the connection prompt. You should be able to connect.
    
    Here ubuntu@ip-10-0-10-166 is the head node, 
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-17.png)
    

In this session, 

1. **Setup the Working Directory and Install the dependencies**
    1. Set up the Python Virtual Environment
    
    ```bash
    python3.9 -m venv ray_env
    ```
    
    b. Activate the venv
    
    ```bash
    source ray_env/bin/activate
    ```
    
    c. Clone the Repo of this Lab
    
    ```bash
    git clone <<<<git repo url>>>>>
    ```
    
    All the necessary files, notebooks are in this repo, change the directory to see the files. We’ll interact with them through Jupyter Lab later in this lab.
    
    ```bash
    cd cloned_repo_dir
    ```
    
    d. Install Ray with pip
    
    ```bash
    pip install --upgrade pip
    pip install ray
    pip install ray['default'] # install this ray dashboard is inaccessible
    ```
    

    💡 use the following source to build the nightly dev version,
    
    ```sh
    https://s3-us-west-2.amazonaws.com/ray-wheels/latest/ray-3.0.0.dev0-cp38-cp38-manylinux2014_x86_64.whl
    ray[default]@ https://s3-us-west-2.amazonaws.com/ray-wheels/latest/ray-3.0.0.dev0-cp38-cp38-manylinux2014_x86_64.whl
    (you may need to add cp39 to the url)
    ```
    
    1. **Start the ray Cluster**
        
        ```bash
        ray start --head --dashboard-host='0.0.0.0'
        ```
    
    If prompted, input ‘n’ to not enable usage stats collection.
    
    After installation, the following stats will emerge,
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-18.png)
    
    To check the status of the Ray cluster, run the following command,
    
    ```bash
    ray status
    ```
    
    The cluster will start and the following stats will show on terminal,
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-19.png)
    
    You can check if the dashboard and ray processes are running on the desired ports with the commands shown in the above image. 
    
    1. **Port Forward the Ray Dashboard to your local machine’s port 8265**
    
    Open another terminal window, and go to directory where your key pair is stored.
    
    ```bash
    ssh -i key-pair-poridhi-poc.pem -N -f -L 8265:localhost:8265 ubuntu@your-head-node's-public-ipv4-DNS
    ```
    
    You can find your head node’s public IPv4 DNS address by selecting the head node’s instance back in the instances tab of AWS console.
    
    Now, open another tab on your browser and go to `localhost:8265` , you can see that the Ray dashboard is running there,
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-20.png)
    
    Go to the cluster tab, you’d see that only 1 node is alive and that is your head node. The other 2 nodes is missing, because we haven’t attached them the cluster yet.
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-21.png)
    

Now, let’s connect the worker nodes to the cluster running on head node.

### For Worker Nodes:

Open individual terminal windows for the instances, and perform the following operations individually,

1. **Connect with the other EC2 instances just like the head node**
2. **Setup the Working Directory and Install the dependencies like the head node**
    1. **Do not clone the GitHub repository mentioned in the head node setup**.
    2. Rest of the operations should be as same as the head node

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-22.png)

1. **Connect with the head node**
    
    After installing all the dependencies on these two nodes, it’s time to connect them with the Ray cluster. 
    
    To connect these with the Ray cluster, you have to use the following command on both of these instances,
    
    ```bash
    ray start --address='10.0.1.166:6379'
    ```
    
    Here, you have to replace the value of address in address flag with the head node’s private IP. The port 6379 should be left as it is.
    
    If the instance fails to connect to the head node, try the following command to increase the number of file descriptors,
    
    ```bash
    ulimit -n 65536; ray start --address='10.0.1.166:6379'
    ```
    
    After connecting to the cluster, the following logs will be shown,
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-23.png)

1. **Observe that the nodes are connected to the cluster on dashboard** 
    
    Now, head back to the Ray dashboard running on localhost’s port 8265 on the browser, and you would find that the worker nodes are added to cluster, (you can also check the status of the cluster with the command **`ray status`** on any node)
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-24.png)
    
    Now that our Ray cluster is running, we can submit jobs and perform tasks on it through Ray Core API and Ray AI runtime libraries. Ray AI runtime libraries have integrations of almost all the popular ML libraries like PyTorch, TensorFlow, Sci-Kit Learn, XGBoost and others. It also has built in support for tools like Prometheus, Grafana, MLFlow, Tensorboard etc. 
    
    We can have different types of nodes for different types of workflow of machine learning like a GPU node for certain tasks and CPU nodes for other tasks. This heterogenous and interoperability nature of Ray clusters make them very useful in every kind of ML tasks, from data processing to model serving and monitoring.
    
    The ML workflow for this lab is not necessarily compute heavy where we will deal with a dataset which is only 5 MB in size. Because, this lab focuses on building intuitions and knowledge around creating and deploying the infrastructure needed for distributed and scalable ML workloads, and also for MLOps with the framework Ray. That’s why we are using one of smallest compute instances on AWS, t2.micro which just has 1 GiB of RAM.
    
    We will now dive into the ML workflows, try to build a ML model for it and serve it through API. All the workflow will be supported by Ray and its Core API.
    

## Install and Run Jupyter Lab:

We will use Jupyter Lab for scripting and running the training/inferencng scripts.

Try to open another terminal window for this task as the terminal have to be on while running the jupyter lab

1. **Head back to the head node and change the directory to the cloned repo**
2. **install the requirements.txt file with pip (this can be done with the requirements.txt too)**
3. **Start the jupyter lab** **(installed in the previous command)**
    
    Make sure you are in the right directory (the cloned repository),
    
    ```bash
    jupyter lab
    ```
    
    And, jupyter lab will be started on port 8888,
    
    ![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-25.png)
    

1. **Port Forward the Jupyter Lab to your local machine’s port 8888**

Open another terminal window, and go to directory where your key pair is stored. 

```bash
ssh -i key-pair-poridhi-poc.pem -N -f -L 8888:localhost:8888 ubuntu@your-head-node's-public-ipv4-DNS
```

Go to the URL provided in the terminal session where the Jupyter Lab is running and you will be able to access the Jupyter Lab,

![alt text](https://raw.githubusercontent.com/AhnafNabil/poridhi.io.intern/main/MLOps%20Lab/Lab%2001/images/image-26.png)

All the necessary files for our problem are already in this working directory.