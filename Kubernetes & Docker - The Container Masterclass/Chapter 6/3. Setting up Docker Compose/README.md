# Installing Docker Compose

## Introduction

In this demo, we will install Docker Compose on our host machine. Docker Compose is a tool that allows us to define and manage multi-container Docker applications. We will fetch the Docker Compose binaries from its official GitHub release, store the binary in the appropriate directory, and make it executable. Finally, we will verify the installation.

## Steps to Install Docker Compose

### Step 1: Fetch Docker Compose Binaries

We will use the `curl` utility to download the Docker Compose binary from its official GitHub release. The binary will be stored in the `/usr/local/bin` directory on our host machine.

1. Open a terminal on your host machine.
2. Run the following command to download the Docker Compose binary:

```bash
sudo curl -L "https://github.com/docker/compose/releases/download/1.22.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
```

This command fetches the Docker Compose binary for the current operating system and architecture, and saves it to `/usr/local/bin/docker-compose`.

### Step 2: Make the Binary Executable

After downloading the binary, we need to make it executable. Use the following command:

```bash
sudo chmod +x /usr/local/bin/docker-compose
```

This command changes the permissions of the Docker Compose binary to make it executable.

### Step 3: Verify the Installation

To ensure Docker Compose is installed correctly, we will check its version. Run the following command:

```bash
docker-compose --version
```

If the installation is successful, you will see output similar to the following:

```bash
docker-compose version 1.22.0, build f46880f
```

This output confirms that Docker Compose version 1.22.0 is installed on our host machine.

## Conclusion

By following these steps, we have successfully installed Docker Compose on our host machine. We fetched the binary using `curl`, made it executable, and verified the installation. Docker Compose is now ready to use for managing multi-container Docker applications.