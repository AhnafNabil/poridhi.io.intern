# Nginx Layer 7 weighted Load Balancing on Node.js servers in AWS

This document outlines the process of setting up a layer 7 load-balanced Node.js application environment using Nginx. The setup consists of two identical Node.js applications, an Nginx server for load balancing. Here we will deploy it in AWS.

<img src="https://github.com/Minhaz00/NodeJS-MySQL/blob/main/10.%20Nginx%20L4%20LB%20NodeJS%20service%20in%20AWS/image/nginxlb-02.PNG?raw=true" />

## Task
Create a load-balanced environment with two `Node.js` applications, `Nginx` as a layer 7 load balancer with weighting parameter and health checkup, and a `MySQL` database, all running in AWS EC2 instance.

## Steps

### Setup AWS: Create VPC, subnets, route table and gateways 

At first, we need to create a VPC in AWS, configure subnet, route tables and gateway.

1. Create a VPC named `my-vpc`
2. Create 2 subnets `public-subnet` and `private-subnet`.
3. Create a public route table named `rt-public` and associate it with `public-subnet`.
4. Create a private route table named `rt-private` and associate it with `private-subnet`
5. Create an internet gateway named `igw` and attach it to `my-vpc`.
6. Create a NAT gateway named `nat-gw` and associate it with `public-subnet`.
7. Configure the route tables to use the internet gateway and NAT gateway.

Here is the `resource-map` of our VPC:

<!-- <img src="https://github.com/Minhaz00/NodeJS-MySQL/blob/main/10.%20Nginx%20L4%20LB%20NodeJS%20service%20in%20AWS/image/image.jpg?raw=true" /> -->

![alt text](https://github.com/Konami33/poridhi.io.intern/blob/main/NGINX%20as%20a%20Load-Balancer/6/images/image.png?raw=true)

### Create EC2 instance

We need to create `3 instances` in EC2. 

#### Create the NodeJS App EC2 Instances:
- Launch two EC2 instances (let's call them `node-app-1` and `node-app-2`) in our private subnet.
- Choose an appropriate AMI (e.g., Ubuntu).
- Configure the instances with necessary security group rules to allow HTTP/HTTPS traffic (typically port 80/443).
- Assign a key pair for SSH access.

#### Create the NGINX EC2 Instance:
- Launch another EC2 instance for the NGINX load balancer (let's call it `nginx-lb`) in our public subnet.
- Configure the instance with a security group to allow incoming traffic on the load balancer port (typically port 80/443) and outgoing traffic to the Flask servers.
- Assign a key pair for SSH access.

### Access the Public Instance via SSH

1. **Set File Permissions**:
   - **For Linux**:
     ```sh
     chmod 400 MyKeyPair.pem
     ```
2. **SSH into the Public Instance**:
   - Open a terminal and run:
     ```sh
     ssh -i MyKeyPair.pem ubuntu@<public_instance_ip>
     ```
   - Replace `<public_instance_ip>` with the public IP address of the public instance.

### Copy the Key Pair to the Public Instance

3. **Copy the Key Pair to the Public Instance**:
   - On your local machine, run the following command to copy the key pair to the public instance:
     ```sh
     scp -i MyKeyPair.pem MyKeyPair.pem ubuntu@<public_instance_ip>:~
     ```
   - Replace `<public_instance_ip>` with the public IP address of the public instance.

### SSH from the Public Instance to the Private Instance

3. **SSH into the Private Instance from the Public Instance**:
   - On the public instance, change the permissions of the copied key pair:
     ```sh
     chmod 400 MyKeyPair.pem
     ```
   - Then, SSH into the private instance:
     ```sh
     ssh -i MyKeyPair.pem ubuntu@<private_instance_ip>
     ```
   - Replace `<private_instance_ip>` with the private IP address of the private instance.


### Install and setup NGINX:

Nginx is very easy to install if we install it from a package manager like apt on Ubuntu or yum in CentOS. It is good for a general proposed load balancer, reverse proxy, and web server. But sometimes we need additional modules to add more function to Nginx that is not included in default installation from the package manager. If that is the case, you need to install the Nginx from source. In this tutorial, we will guide you step by step on how to install Nginx from source on ubuntu.

#### Sudo Privileges
Before starting, make sure that we have no permission issue on the installation & configuration.

```bash
sudo su
```

#### Install Dependencies
Run this command to install Nginx dependencies

```bash
apt update -y && apt-get install git build-essential libpcre3 libpcre3-dev zlib1g zlib1g-dev libssl-dev libgd-dev libxml2 libxml2-dev uuid-dev
```

#### Download Nginx Source Code
Before you download the Nginx source code, you can visit http://nginx.org/en/download.html to see the Nginx version available now. After that you can download them by running this command:

```bash
wget http://nginx.org/download/nginx-<version>.tar.gz
```
Now, the latest stable version is 1.26.1, so for me, I will download the nginx-1.26.1 version

```bash
wget http://nginx.org/download/nginx-1.26.1.tar.gz
```

Extract the downloaded file

```bash
tar -zxvf nginx-1.26.1.tar.gz
```

#### Build & Install Nginx
After extract the file, go to the nginx directory

```bash
cd nginx-1.20.1
```
Now is the time to configure Nginx that suits your need, this is where you put in the module you want to include in Nginx using the ./configure command. The full documentation is in here: Building Nginx from Sources. For now, I will give you the minimum configure option so you can build a good load balancer, reverse proxy, or webserver. Run this command to configure Nginx:

```bash
./configure \
    --prefix=/etc/nginx \
    --conf-path=/etc/nginx/nginx.conf \
    --error-log-path=/var/log/nginx/error.log \
    --http-log-path=/var/log/nginx/access.log \
    --pid-path=/run/nginx.pid \
    --sbin-path=/usr/sbin/nginx \
    --with-http_ssl_module \
    --with-http_v2_module \
    --with-http_stub_status_module \
    --with-http_realip_module \
    --with-file-aio \
    --with-threads \
    --with-stream \
    --with-stream_ssl_preread_module
```

After that, run this command to build & install the Nginx

```bash
make && make install
```

To verify the installation, you can check the Nginx version

```bash
nginx -V
```

### Configure NGINX as an L4 Load Balancer

First, move to work directory to the Nginx configuration folder
```sh
cd /etc/nginx
```
Backup the default Nginx configuration file
```sh
mv nginx.conf nginx.conf.old
```
Open the NGINX configuration file.
```sh
sudo vim nginx.conf
```
**nginx.conf:**

Add the following configuration to the file:

```sh
http {
    upstream nodejs_backend {
        server <Node-app-1-ip>:3001 weight=3;  # Node-app-1 with weight 3
        server <Node-app-2-ip>:3002 weight=1;  # Node-app-2 with weight 1
    }

    server {
        listen 80 default_server;
        listen [::]:80 default_server;
        server_name _;

        location / {
            proxy_pass http://nodejs_backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_cache_bypass $http_upgrade;

            # Additional headers for proper proxying
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Logging
        log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                        '$status $body_bytes_sent "$http_referer" '
                        '"$http_user_agent" "$http_x_forwarded_for"';

        access_log /var/log/nginx/access.log main; # Access log path
        error_log /var/log/nginx/error.log; # Error log path
    }
}
```

### Explanation and Adjustments

1. **`upstream` Block:**
   - Defines the backend servers `nodejs_backend` with their respective weights (`weight=3` for Node-app-1 and `weight=1` for Node-app-2), which is correct for load balancing.

2. **`server` Block:**
   - Listens on port 80 (`listen 80`) for incoming HTTP requests.
   - Uses `proxy_pass` to forward requests to the defined upstream servers (`http://nodejs_backend`).

3. **Headers and Proxy Settings:**
   - `proxy_http_version 1.1` ensures HTTP/1.1 is used for proxying.
   - `proxy_set_header` directives are used to forward headers such as `Upgrade`, `Connection`, `Host`, `X-Real-IP`, `X-Forwarded-For`, and `X-Forwarded-Proto` for proper proxying and client IP preservation.

4. **Logging:**
   - `log_format main` defines the format for access logs (`access_log`) stored in `/var/log/nginx/access.log`.
   - Error logs (`error_log`) are directed to `/var/log/nginx/error.log`.

### Recommendations

1. **Configuration Test:**
   - Before deploying changes, always test the configuration file for syntax errors:

     ```bash
     sudo nginx -t
     ```

2. **Reload Nginx:**
   - If the configuration test is successful, reload Nginx to apply changes:

     ```bash
     sudo systemctl reload nginx
     ```

3. **Monitor Logs:**
   - Monitor `/var/log/nginx/access.log` and `/var/log/nginx/error.log` to troubleshoot and analyze HTTP requests and potential errors.

### Setup Node.js Applications on Different Instances

1. **SSH into each Node.js instance and perform the following steps:**

2. **Update the package list and install npm:**
   ```bash
   sudo apt update
   sudo apt install npm
   ```

3. **Create a directory for your Node.js application and navigate into it:**
   ```bash
   mkdir node-app
   cd node-app
   ```

4. **Initialize a new Node.js project:**
   ```bash
   npm init -y
   ```

5. **Install Express.js:**
   ```bash
   npm install express
   ```

6. **Create the application file:**
   ```bash
   nano index.js
   ```

7. **Add the following code to `index.js`:**
   ```javascript
   const express = require('express');
   const app = express();
   const port = process.env.PORT;

   app.get('/', (req, res) => {
     res.status(200).send(`Hello, from Node App on PORT: ${port}!`);
   });

   app.listen(port, () => {
     console.log(`App running on http://localhost:${port}`);
   });
   ```

8. **Set the port environment variable and start the application:**

For the first instance:
   ```bash
   export PORT=3001  # or 3002 for the second instance
   node index.js
   ```
For the second instance":
   ```bash
   export PORT=3001  # or 3002 for the second instance
   node index.js
   ```

### Verify the Load Balancing Setup

1. **Access the public IP of your NGINX load balancer in a web browser**:
   ```http
   http://<load-balancer-public-ip>
   ```

   ![alt text](https://github.com/Konami33/poridhi.io.intern/blob/main/NGINX%20as%20a%20Load-Balancer/6/images/Screenshot%202024-06-26%20214643.png?raw=true)

   ![alt text](https://github.com/Konami33/poridhi.io.intern/blob/main/NGINX%20as%20a%20Load-Balancer/6/images/Screenshot%202024-06-26%20214707.png?raw=true)

   <!-- <img src="https://github.com/Minhaz00/NodeJS-MySQL/blob/main/10.%20Nginx%20L4%20LB%20NodeJS%20service%20in%20AWS/image/image2.jpg?raw=true" />
   
    <img src="https://github.com/Minhaz00/NodeJS-MySQL/blob/main/10.%20Nginx%20L4%20LB%20NodeJS%20service%20in%20AWS/image/image1.jpg?raw=true" /> -->


By following these steps, you set up a Layer-7 load balancer using NGINX to distribute traffic between two Node.js applications running on different instances. This ensures that the load is balanced and provides high availability for your applications.