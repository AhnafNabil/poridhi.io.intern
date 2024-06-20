# In-Depth Container Inspection and Committing Changes

In this demo, we will delve deeper into understanding our containers by using the `docker inspect` and `docker commit` commands. These commands help us gather detailed information about our containers and commit changes to create new images.

## Listing Containers

First, let's list the current containers to see which ones are available.

### Command:
```sh
docker ps -a
```
- **Explanation**: This command lists all containers, including those that are running, stopped, or exited.

## Inspecting a Container

We will use the `docker inspect` command to get detailed information about the `my-ubuntu` container.

### Steps:

1. **Inspect the Container:**
    - **Command**: `docker inspect my-ubuntu`
    - **Explanation**: This command provides a detailed JSON description of the container.

    ```sh
    docker inspect my-ubuntu
    ```

    - **Output**: The output includes a plethora of information, such as container ID, creation timestamp, state, image details, network settings, and more.

2. **Interpreting Key Information:**
    - **Container ID**: A unique identifier for the container.
    - **Creation Timestamp**: The time when the container was created.
    - **State**: Indicates the current state of the container (running, paused, restarting, dead).
    - **Process ID**: The process ID of the container on the host (e.g., 694).
    - **Paths**: Host path, log path, and configuration path.
    - **Network Settings**: Details about the container's network configuration, including the IP address.

3. **Filtering Specific Information:**
    - **Command**: `docker inspect --format '{{ .NetworkSettings.IPAddress }}' my-ubuntu`
    - **Explanation**: This command uses the `--format` flag to narrow down the output to just the IP address of the container.

    ```sh
    docker inspect --format '{{ .NetworkSettings.IPAddress }}' my-ubuntu
    ```

    - **Output**: The IP address of the `my-ubuntu` container.

## Committing Changes to a Container

To effectively use the `docker commit` command, we need to make at least one change to the container's state after it is created from the image.

### Steps:

1. **Start the Container:**
    - **Command**: `docker container start my-ubuntu`
    - **Explanation**: This command starts the `my-ubuntu` container if it's not already running.

    ```sh
    docker container start my-ubuntu
    ```

2. **Execute a Command in the Container:**
    - **Command**: `docker exec -it my-ubuntu bash`
    - **Explanation**: This command opens an interactive terminal inside the container.

    ```sh
    docker exec -it my-ubuntu bash
    ```

3. **Make a Change:**
    - **Command**: `apt-get update`
    - **Explanation**: This command updates the package lists inside the container, altering its state.

    ```sh
    apt-get update
    ```

4. **Exit the Container:**
    - **Command**: `exit`
    - **Explanation**: This command exits the container's terminal.

    ```sh
    exit
    ```

5. **Commit the Changes:**
    - **Command**: `docker commit my-ubuntu <your-dockerhub-username>/updated-ubuntu:1.0`
    - **Explanation**: This command commits the changes to the container, creating a new image with the specified tag.

    ```sh
    docker commit my-ubuntu <your-dockerhub-username>/updated-ubuntu:1.0
    ```

    - **Note**: Ensure you are logged into your Docker Hub account to push the changes.

6. **List the Images:**
    - **Command**: `docker images`
    - **Explanation**: This command lists all the Docker images available locally, including the newly committed image.

    ```sh
    docker images
    ```

## Summary

In this session, we covered:

- **Inspecting Containers**:
    - Command: `docker inspect <container_name>`
    - Example: `docker inspect my-ubuntu`
    - Outcome: Provides detailed JSON output about the container.

- **Filtering Specific Information**:
    - Command: `docker inspect --format '{{ .NetworkSettings.IPAddress }}' <container_name>`
    - Example: `docker inspect --format '{{ .NetworkSettings.IPAddress }}' my-ubuntu`
    - Outcome: Outputs the IP address of the container.

- **Committing Changes to Containers**:
    - Command: `docker commit <container_name> <dockerhub-username>/<new-image>:<tag>`
    - Example: `docker commit my-ubuntu <your-dockerhub-username>/updated-ubuntu:1.0`
    - Outcome: Creates a new image with the committed changes.

By following these steps, you can inspect containers in detail and commit changes to create new Docker images. In the next lecture, we will explore port mapping and how it allows us to expose container services to the host.