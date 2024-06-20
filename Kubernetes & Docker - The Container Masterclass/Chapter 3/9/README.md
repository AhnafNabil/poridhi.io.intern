# Attaching and Executing Commands in Containers

In this demo, we will explore the `docker container attach` and `docker exec` commands. These commands allow us to interact with running containers in different ways.

## Listing Containers

First, let's list the current containers to see which ones are available.

### Command:
```sh
docker ps -a
```
- **Explanation**: This command lists all containers, including those that are running, stopped, or exited.

## Attaching to a Container

Next, we will use the `docker container attach` command to attach to the `my-busybox` container.

### Steps:

1. **Attach to the Container:**
    - **Command**: `docker container attach my-busybox`
    - **Explanation**: This command attaches the standard input, output, and error of the `my-busybox` container to the terminal of our Docker client.

    ```sh
    docker container attach my-busybox
    ```

2. **Interacting with the Container:**
    - Once attached, we can interact with the container's terminal. For example, we can list the directories in the container's root environment.
    - **Command**: `ls`
    - **Explanation**: This command lists the directories in the container's root environment.

    ```sh
    ls
    ```

3. **Exit the Container:**
    - **Command**: `exit`
    - **Explanation**: This command exits the container's terminal, returning us to the host's terminal.

    ```sh
    exit
    ```

4. **Verify Container Status:**
    - **Command**: `docker ps -a`
    - **Explanation**: This command lists all containers again to verify their status.

    ```sh
    docker ps -a
    ```

    - **Outcome**: The `my-busybox` container will be in the exited state because exiting the attachment also stops the container.

## Executing Commands in a Container

Now, let's use the `docker exec` command to run commands in the container without stopping it.

### Steps:

1. **Start the Container:**
    - **Command**: `docker container start my-busybox`
    - **Explanation**: This command starts the `my-busybox` container again.

    ```sh
    docker container start my-busybox
    ```

2. **Execute a Command in the Container:**
    - **Command**: `docker exec -it my-busybox pwd`
    - **Explanation**: This command executes the `pwd` (print working directory) command in the `my-busybox` container with interactive and TTY enabled.

    ```sh
    docker exec -it my-busybox pwd
    ```

    - **Outcome**: The command outputs `/`, indicating the root directory of the `my-busybox` container.

3. **Verify Container Status:**
    - **Command**: `docker ps -a`
    - **Explanation**: This command lists all containers again to verify their status.

    ```sh
    docker ps -a
    ```

    - **Outcome**: Unlike `attach`, using `exec` does not stop the container, so `my-busybox` will still be running.

## Summary

In this session, we covered:

- **Attaching to a Container**:
    - Command: `docker container attach <container_name>`
    - Example: `docker container attach my-busybox`
    - Outcome: Interact directly with the container's terminal. Exiting stops the container.

- **Executing Commands in a Container**:
    - Command: `docker exec -it <container_name> <command>`
    - Example: `docker exec -it my-busybox pwd`
    - Outcome: Executes a command in the container without stopping it.

By following these steps, you can effectively interact with and manage running Docker containers. In the next lecture, we will explore more application-related operations with our containers.