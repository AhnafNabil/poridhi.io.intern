# Port Mapping in Docker Containers

In this demo, we will learn how to map our host machine's ports to container ports. This is crucial for exposing services running inside Docker containers to the outside world.

## Mapping Host Port to Container Port

### Step-by-Step Instructions:

1. **Run a Container with Port Mapping:**
    - **Command**: `docker run -d -p 8080:80 <image_name>`
    - **Explanation**: This command runs a container in detached mode (`-d`), mapping port 8080 on the host to port 80 in the container. The `-p` flag is used for port mapping.

    ```sh
    docker run -d -p 8080:80 <your_image_name>
    ```

    - **Example**: If your image is `nginx`, the command would be:

    ```sh
    docker run -d -p 8080:80 nginx
    ```

    - **Output**: The container ID will be displayed, and you can check the running container with the `docker ps` command.

    ```sh
    docker ps
    ```

    - **Expected Output**: The output will include a column `PORTS` showing `0.0.0.0:8080->80/tcp`.

2. **Run Another Container with Automatic Port Mapping:**
    - **Command**: `docker run -d -P --name cont_nginx-A <image_name>`
    - **Explanation**: This command runs a container with automatic port mapping (`-P`), allowing Docker to map container ports to available host ports. The `--name` flag is used to name the container.

    ```sh
    docker run -d -P --name cont_nginx-A <your_image_name>
    ```

    - **Example**: If your image is `nginx`, the command would be:

    ```sh
    docker run -d -P --name cont_nginx-A nginx
    ```

    - **Output**: The container ID will be displayed, and you can check the running container with the `docker ps` command.

    ```sh
    docker ps
    ```

    - **Expected Output**: The output will include a column `PORTS` showing mappings like `0.0.0.0:32768->80/tcp`.

3. **Check Port Mapping Information:**
    - **Command**: `docker container port <container_name>`
    - **Explanation**: This command displays the port mappings for a specific container.

    ```sh
    docker container port cont_nginx-A
    ```

    - **Expected Output**: It will show the mapping, such as `80/tcp -> 0.0.0.0:32768`.

4. **Verify Port Mapping in a Web Browser:**
    - **Localhost URL**: `http://localhost:8080`
    - **Explanation**: Open a web browser and navigate to `http://localhost:8080` to see the Nginx homepage, indicating that port mapping was successful.

    ```sh
    # Open a web browser and go to:
    http://localhost:8080
    ```

    - **Expected Result**: The Nginx homepage should be displayed, confirming that the port mapping works.

    - **Repeat for the Other Container**: Similarly, navigate to the mapped port of the second container (e.g., `http://localhost:32768`).

## Summary

In this session, we covered how to map host ports to container ports, both manually and automatically, using the following steps:

1. **Manual Port Mapping**: 
    - Command: `docker run -d -p 8080:80 <image_name>`
    - Example: `docker run -d -p 8080:80 nginx`

2. **Automatic Port Mapping**:
    - Command: `docker run -d -P --name cont_nginx-A <image_name>`
    - Example: `docker run -d -P --name cont_nginx-A nginx`

3. **Check Port Mapping**:
    - Command: `docker container port <container_name>`
    - Example: `docker container port cont_nginx-A`

4. **Verify in Browser**:
    - URL: `http://localhost:8080`
    - Expected: Nginx homepage should be displayed.

By following these steps, you can effectively map host machine ports to your container ports, making your containerized applications accessible from the host machine. In the next lecture, we'll focus on cleaning up our workspace.