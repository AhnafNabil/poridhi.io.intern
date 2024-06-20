# Advanced Dockerfile: Configuration Instructions

In this demo, we will further explore writing Dockerfiles by focusing on configuration instructions. These instructions allow us to set up and configure our Docker images in a more controlled and efficient manner.

## Step 1: Verify Current Directory and Navigate

Ensure we are in the `CC_Docker` directory, which contains individual directories for each demo. Navigate to the `D-2` directory after creating it.

```sh
cd D-2
```

## Step 2: Open the Dockerfile

In the `D-2` directory, Create an empty Dockerfile using the `touch` command. Open it using a text editor such as `nano`.

```sh
touch Dockerfile
```

```sh
nano Dockerfile
```

## Step 3: Examine the Dockerfile

The Dockerfile includes the following instructions:

1. **FROM Instruction:** Specifies the base image. Here, Ubuntu 16.04 is used.

    ```Dockerfile
    FROM ubuntu:16.04
    ```

2. **RUN Instructions:** Execute commands on top of the base image, with each command creating a new layer.

    ```Dockerfile
    RUN apt-get update -y && \
        apt-get install -y curl && \
        apt-get clean

    RUN mkdir -p /home/Codes
    ```

3. **ENV Instructions:** Set environment variables within the image.

    ```Dockerfile
    ENV USER nabil
    ENV SHELL /bin/bash
    ENV LOGNAME nabil
    ```

4. **CMD Instruction:** Specifies the command to run within the container (covered in later demos).

    ```Dockerfile
    CMD ["bash"]
    ```

## Step 4: Build the Docker Image

Use the `docker build` command to build the Dockerfile into an image. Tag the image as `img_run-env` for easy identification.

```sh
docker build -t img_run-env .
```

## Step 5: Verify the Build Process

The build process involves several steps:

1. **Setting up the base image:** Directly sets the base image without using `ARG`.
2. **Executing RUN instructions:** Performs the commands to update the OS, install `curl`, and create the `/home/Codes` directory.
3. **Setting ENV variables:** Sets the `USER`, `SHELL`, and `LOGNAME` environment variables.

Once the build is complete, verify the created image using the `docker images` command.

```sh
docker images
```

The `img_run-env` image should be listed among the available images.

## Step 6: Run the Docker Image as a Container

Run the image as a container using the `docker run -itd` command. The `-itd` flags indicate interactive, teletype-enabled, and detached modes, respectively. Name the container `cont_run-env`.

```sh
docker run -itd --name cont_run-env img_run-env
```

The command returns a unique container ID, indicating the container is running.

## Step 7: Verify the Running Container

List the running containers to ensure `cont_run-env` is active.

```sh
docker ps
```

## Step 8: Access the Running Container

Access the container's bash shell using the `docker exec` command.

```sh
docker exec -it cont_run-env bash
```

## Step 9: Verify Environment Variables and Directory

Within the container, verify the environment variables and directory created by the Dockerfile.

1. **Check Environment Variables:**

    ```sh
    echo $USER
    echo $SHELL
    echo $LOGNAME
    ```

    The values should match those set in the Dockerfile (`USER=nabil`, `SHELL=/bin/bash`, `LOGNAME=nabil`).

2. **Check Directory Structure:**

    Navigate to the home directory and list its contents to verify the creation of the `/home/Codes` directory.

    ```sh
    cd /home
    ls
    ```

## Step 10: Exit the Container

Exit the container to return to the host environment.

```sh
exit
```

## Summary

- **FROM Instruction:** Sets the base image.
- **RUN Instruction:** Executes commands, creating separate layers.
- **ENV Instruction:** Sets environment variables.
- **CMD Instruction:** Specifies the default command for the container.

By following these steps, we have created a Dockerfile with configuration instructions, built a Docker image, and verified the setup by running a container. This demonstrates the power and flexibility of Dockerfile instructions in configuring and managing Docker images.