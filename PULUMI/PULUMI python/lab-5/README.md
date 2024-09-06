# Create AWS Infrastructure with GitHub Actions and SSH Access

In this lab, we will set up a Virtual Private Cloud (VPC) in AWS with both public and private subnets, launch EC2 instances in each subnet, and establish secure communication between these instances. Specifically, we will:

1. Configure Pulumi.
2. Set up a VPC with a public and a private subnet.
3. Create an Internet Gateway (IGW) for the public subnet.
4. Create a NAT Gateway for the private subnet.
5. Launch an EC2 instance in the public subnet with a public IP address.
6. Launch EC2 instances in the private subnet.
7. Attach Github action to create the resources using PULUMI. Basically Github action will trigger the pulumi to create the resouces.
8. Establish SSH access to the EC2 instances.

![alt text](./images/github-action.drawio.svg)

## Step by step

## Step 01: Set Up a Pulumi Project
- Create a new directory for your project and navigate into it:

  ```sh
  mkdir aws-infra
  cd aws-infra
  ```

- Install Python venv

  ```sh
  sudo apt update
  sudo apt install python3.8-venv
  ```
- Run the following command to create a new Pulumi project:

    ```sh
    pulumi new aws-python
    ```
    Follow the prompts to set up your project. It will create the necessary folders for pulumi setup.

- Create SSH Key-pair file

  ```sh
  aws ec2 create-key-pair --key-name MyKeyPair --query 'KeyMaterial' --output text > MyKeyPair.pem
  ```

- Set the permission of the key file

  ```sh
  chmod 400 MyKeyPair.pem
  ```

- Update the `__main__.py` file according to this:

```python

```

## Step 02: Create GitHub Actions Workflow

- Create a GitHub Actions workflow in this directory `.github/workflows/deploy.yml` of your project.

  ```yaml
  name: Deploy Infrastructure

  on:
    push:
      branches:
        - main

  jobs:
    deploy:
      runs-on: ubuntu-latest

      steps:
        - name: Checkout code
          uses: actions/checkout@v3

        - name: Set up Python
          uses: actions/setup-python@v2
          with:
            python-version: '3.x'

        - name: Install dependencies
          run: |
            python -m venv venv
            source venv/bin/activate
            pip install pulumi pulumi-aws

        - name: Configure AWS credentials
          uses: aws-actions/configure-aws-credentials@v2
          with:
            aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
            aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
            aws-region: ap-southeast-1

        # login
        - name: Pulumi login
          env:
            PULUMI_ACCESS_TOKEN: ${{ secrets.PULUMI_ACCESS_TOKEN }}
          run: pulumi login

        # select stack
        - name: Pulumi stack select
          run: pulumi stack select dev

        - name: Set public key environment variable
          run: echo "PUBLIC_KEY=${{ secrets.PUBLIC_KEY }}" >> $GITHUB_ENV

        # refresh
        - name: Pulumi refresh
          run: pulumi refresh --yes

        # up
        - name: Pulumi up
          run: pulumi up --yes

        # save outputs
        - name: Save Pulumi outputs
          id: pulumi_outputs
          run: |
            MASTER_IP=$(pulumi stack output master_public_ip)
            WORKER1_IP=$(pulumi stack output worker1_public_ip)
            WORKER2_IP=$(pulumi stack output worker2_public_ip)
            NGINX_IP=$(pulumi stack output nginx_instance)

            echo "MASTER_IP=$MASTER_IP" >> $GITHUB_ENV
            echo "WORKER1_IP=$WORKER1_IP" >> $GITHUB_ENV
            echo "WORKER2_IP=$WORKER2_IP" >> $GITHUB_ENV
            echo "NGINX_IP=$NGINX_IP" >> $GITHUB_ENV

          env:
            PULUMI_ACCESS_TOKEN: ${{ secrets.PULUMI_ACCESS_TOKEN }}


        - name: Debug Environment Variables
          run: |
            echo "MASTER_IP=${{ env.MASTER_IP }}"
            echo "WORKER1_IP=${{ env.WORKER1_IP }}"
            echo "WORKER2_IP=${{ env.WORKER2_IP }}"
            echo "NGINX_IP=${{ env.NGINX_IP }}"
  ```

- Commit your changes and push to your GitHub repository.

  ```sh
  git add .
  git commit -m "Add Pulumi infrastructure and workflow"
  git push origin main
  ```

## Step 03: SSH connectivity

In our architecture we have a public instance and private instances. To connect to these instances, we need to first SSH into our public instance and then from there we can SSH into the private instances.

- **First copy the key file into to the public instance `~.ssh/` directory:**

  ```sh
  scp -i MyKeyPair.pem MyKeyPair.pem ubuntu@public-instance-ip:~/.ssh/
  ```
- Now SSH into the public instance

  ```sh
  ssh -i MyKeyPair.pem ubuntu@public-instance-ip
  ```
- **Then SSH into the private instances from the public instance:**

  ```sh
  ssh -i MyKeyPair.pem ubuntu@private-instance-ip
  ```

### Conclusion

So, we have setup pulumi project, created github action to automatically create the aws resources using pulumi. Then we have also SSH into the instances.