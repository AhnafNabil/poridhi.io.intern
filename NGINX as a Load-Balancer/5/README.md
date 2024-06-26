# Nginx Layer 4 weighted Load Balancing on NodeJS-MySQL App in AWS

This document outlines the process of setting up a `layer 4` load-balanced Node.js application environment using Nginx. The setup consists of two identical Node.js applications, an `Nginx` server for weighted load balancing and health checkup. We will use `mysql` database for this app. We will also deploy it in AWS.

<img src="https://github.com/Minhaz00/NodeJS-MySQL/blob/main/11.%20Nginx%20L4%20LB%20NodeJS-MySQL%20App%20in%20AWS/images/nginx-01.PNG?raw=true" />

## Task
Create a load-balanced environment with two `Node.js` applications, `Nginx` as a load balancer with weighting parameter and health checkup, and a `MySQL` database, all running in AWS EC2 instance.


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

<img src="https://github.com/Minhaz00/NodeJS-MySQL/blob/main/10.%20Nginx%20L4%20LB%20NodeJS%20service%20in%20AWS/image/image.jpg?raw=true" />

### Create and setup EC2 instances

We need to create `3 instances` in EC2. One in the public subnet for nginx server, second one is for Nodejs Application and third one is mysql. Both will be in the private subnet.

#### Create the NodeJS App EC2 Instances:
- Launch two EC2 instances (let's call them `node-app-1` and `node-app-2`) in our private subnet.
- Choose an appropriate AMI (e.g., Ubuntu).
- Configure the instances with necessary security group rules to allow HTTP/HTTPS traffic (typically port 80/443).
- Assign a key pair for SSH access.

#### Create the NGINX EC2 Instance:
- Launch another EC2 instance for the NGINX load balancer (let's call it `nginx-lb`) in our public subnet.
- Configure the instance with a security group to allow incoming traffic on the load balancer port (typically port 80/443) and outgoing traffic to the NodeJS servers.
- Assign a key pair for SSH access.

#### Create the mysql EC2 Instance:
- Launch another EC2 instance for the MySQL Database (let's call it `mysql`).


### Access the Public Instance via SSH

1. *Set File Permissions*:
   - *For Windows*: Open PowerShell and navigate to the directory where MyKeyPair.pem is located. Then, use the following command to set the correct permissions:
     ```powershell
     icacls MyKeyPair.pem /inheritance:r
     icacls MyKeyPair.pem /grant:r "$($env:USERNAME):(R)"
     ```

   - *For Linux*:
     ```sh
     chmod 400 MyKeyPair.pem
     ```


2. *SSH into the Public Instance*:
   - Open a terminal and run:
     ```sh
     ssh -i MyKeyPair.pem ubuntu@<public_instance_ip>
     ```
   - Replace <public_instance_ip> with the public IP address of the public instance.

#### Copy the Key Pair to the Public Instance

3. *Copy the Key Pair to the Public Instance*:
   - On your local machine, run the following command to copy the key pair to the public instance:
     ```sh
     scp -i MyKeyPair.pem MyKeyPair.pem ubuntu@<public_instance_ip>:~
     ```
   - Replace <public_instance_ip> with the public IP address of the public instance.

#### SSH from the Public Instance to the Private Instance

3. *SSH into the Private Instance from the Public Instance*:
   - On the public instance, change the permissions of the copied key pair:
     ```sh
     chmod 400 MyKeyPair.pem
     ```
   - Then, SSH into the private instance:
     ```sh
     ssh -i MyKeyPair.pem ubuntu@<private_instance_ip>
     ```
   - Replace <private_instance_ip> with the private IP address of the private instance.

## Set up MySQL

Connect to `mysql` instance install `Docker` and run the following command to run mysql database container:

```bash
docker run --name my_mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=my_db -e MYSQL_USER=my_user -e MYSQL_PASSWORD=my_password -p 3306:3306 -v mysql_data:/var/lib/mysql -d mysql:latest
```

This command creates and runs a MySQL database container with specified database, user, and password and volume mounting for persistent storage, accessible externally on port 3306.

You need to set inbound security-group to give access on port `3306`.



## Set up Node.js Applications

Connect to `node-app-1` and configure as follows:

### Create Node App 1

Install npm:
```sh
sudo apt update
sudo apt install npm
```

Then, 
```bash
mkdir Node-app
cd Node-app
npm init -y
npm install express mysql2 body-parser
```

Create `index.js` in the Node-app-1 directory with the following code:

```javascript
const express = require('express');
const mysql = require('mysql2');
const bodyParser = require('body-parser');

const app = express();
const port = process.env.PORT;

const dbConfig = {
  host: "<mysql-instance-private-ip>",
  user: 'my_user',
  password: 'my_password',
  database: 'my_db'
};

// Middleware to parse JSON bodies
app.use(bodyParser.json());

function createConnection() {
  const connection = mysql.createConnection(dbConfig);
  connection.connect(error => {
    if (error) {
      console.error('Error connecting to the database:', error);
      return null;
    }
    console.log('Connected to MySQL database');
  });
  return connection;
}

function createTable() {
  const connection = createConnection();
  if (connection) {
    const createTableQuery = `
      CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE
      )
    `;
    connection.query(createTableQuery, (error, results) => {
      connection.end();
      if (error) {
        console.error('Error creating table:', error);
      } else {
        console.log('Table "users" ensured to exist');
      }
    });
  }
}

// Ensure the table is created when the server starts
createTable();

app.get('/', (req, res) => {
  
  const connection = createConnection();
  if (connection) {
    res.status(200).json({ message: `Hello, from Node App on PORT: ${port} ! Connected to MySQL database.` });
    connection.end();
  } else {
    res.status(500).json({ message: 'Failed to connect to MySQL database' });
  }
});


// GET route to fetch all users
app.get('/users', (req, res) => {
  const connection = createConnection();
  if (connection) {
    connection.query('SELECT * FROM users', (error, results) => {
      connection.end();
      if (error) {
        return res.status(500).json({ message: 'Error fetching users', error });
      }
      res.json(results);
    });
  } else {
    res.status(500).json({ message: 'Failed to connect to MySQL database' });
  }
});

// POST route to add a new user
app.post('/users', (req, res) => {
  const connection = createConnection();
  const { name, email } = req.body;

  if (!name || !email) {
    return res.status(400).json({ message: 'Name and email are required' });
  }

  if (connection) {
    const query = 'INSERT INTO users (name, email) VALUES (?, ?)';
    connection.query(query, [name, email], (error, results) => {
      connection.end();
      if (error) {
        return res.status(500).json({ message: 'Error adding user', error });
      }
      res.status(201).json({ message: 'User added', userId: results.insertId });
    });
  } else {
    res.status(500).json({ message: 'Failed to connect to MySQL database' });
  }
});

app.listen(port, () => {
  console.log(`App running on port :${port}`);
});

```

Replace `<mysql-instance-private-ip>` with your mysql instance private ip.

### Create Node App 2

Connect to the `Node-app-2` instance. 
Here, do the similar steps as `Node-app-1`.
 

### Start Node.js Applications
Navigate to each Node.js application directory and run:
```bash
export PORT=5001
node index.js
```

Do the same for both nodejs app instances. Make sure they are running and connected to the database properly. Note that, you need to set the `PORT=5002` in `Node-app-2` instance.


## Set up Nginx

Now, connect to teh `nginx` instance and create a `nginx.conf` file and a `Dockerfile`. You also need to install `Docker`. 

### Create Nginx Configuration
```bash
mkdir Nginx
cd Nginx
```

Create `nginx.conf` in the Nginx directory with the following configuration:

```sh
events {}

stream {
    upstream nodejs_backend {
        server <Node-app-1 private-ip>:3001 weight=3; # Node-app-1 with weight 3
        server <Node-app-2 private-ip>:3002 weight=1; # Node-app-2 with weight 1

        # Optional: Specify the maximum number of connections per server
        # max_conns 100;
    }

    server {
        listen 80;

        proxy_pass nodejs_backend;

        # Enable TCP load balancing
        proxy_connect_timeout 1s;
        proxy_timeout 3s;

        # Health check configuration
        health_check interval=5s fails=2 passes=2;
        health_check_timeout 2s; # Time to wait for a health check response
        health_check_connect_timeout 1s; # Time to wait for a connection attempt
        health_check_match using=string expected="200 OK"; # Expected response for a successful check

        # Logging
        access_log /var/log/nginx/access.log; # Access log path
        error_log /var/log/nginx/error.log; # Error log path
    }
}

```

Replace the `pubic ip` of nodejs app according to `your ec2 instances`.

**Explanation**:
- `Weight-Based Load Balancing:` 
    Each server in the upstream block has a weight parameter. The higher the weight, the more requests are sent to that server. In this case, <Node-app-1 private-ip>:3001 has a weight of 3, and <Node-app-2 private-ip>:3002 has a weight of 1.

- `Health Checks:`

    - `health_check_timeout 2s;`: Sets the time to wait for a health check response.
    - `health_check_connect_timeout 1s;`: Sets the time to wait for a connection attempt during a health check.
    - `health_check_match using=string expected="200 OK";`: Specifies the expected response for a successful health check. This example assumes the health check endpoint returns a "200 OK" status.


This configuration provides a simple weighted round-robin load balancing with a healthcheckup across our Node.js applications at the TCP level, which can be more efficient than HTTP-level load balancing for certain use cases.

### Create Dockerfile for Nginx
Create a file named `Dockerfile` in the Nginx directory with the following content:

```Dockerfile
FROM nginx:latest
COPY nginx.conf /etc/nginx/nginx.conf
```

### Build Nginx Docker Image
```bash
docker build -t custom-nginx .
```

This command builds a Docker image for Nginx with our custom configuration.

### Run Nginx Container
```bash
docker run -d -p 80:80 --name my_nginx custom-nginx
```

This command starts the Nginx container with our custom configuration.


## Verification

1. Visit `http://<nginx-public-ip>` in a web browser. You should see a response from one of the Node.js applications.

2. To test the user operations, use the following endpoints:
   - GET `http://<nginx-public-ip>/users` (to fetch all users)
   - POST `http://<nginx-public-ip>/users` (to add a new user)

   For the POST request, use a JSON body like:
   ```json
   {
     "name": "John Doe",
     "email": "john@example.com"
   }
   ```

    I am using `Postman` for `POST` method. After adding the user, let's try to see the user using `GET` request: `http://<nginx-public-ip>/users`


3. Refresh the browser or make multiple requests to observe the load balancing in action. You should see responses alternating between Node-app-1 and Node-app-2 running on different ports.

    Example:

    <img src="https://github.com/Minhaz00/NodeJS-MySQL/blob/main/11.%20Nginx%20L4%20LB%20NodeJS-MySQL%20App%20in%20AWS/images/app1.png?raw=true" />

    <img src="https://github.com/Minhaz00/NodeJS-MySQL/blob/main/11.%20Nginx%20L4%20LB%20NodeJS-MySQL%20App%20in%20AWS/images/app2.png?raw=true" />
   
## Conclusion

By configuring Nginx with the stream module as described, you have effectively set up Nginx as a `Layer 4 load balancer`. It handles TCP connections based on IP addresses and port numbers, making routing decisions at the transport layer without inspecting higher-layer protocols like HTTP. This setup is suitable for scenarios where efficient TCP load balancing across multiple backend services is required.