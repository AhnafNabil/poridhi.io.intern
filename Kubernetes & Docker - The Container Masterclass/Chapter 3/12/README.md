# Removing Docker Containers

In this demo, we will explore different methods to remove Docker containers. Over time, containers can accumulate, and it's essential to clean up unnecessary ones to free up resources and maintain a tidy environment.

## Listing All Containers

Before we start removing containers, let's list all of them to see what we have:

```sh
docker ps -a
```

## Removing Containers

### 1. Removing a Single Stopped Container by Name

- **Command**: `docker rm <container_name>`
- **Example**: Let's remove a container named `cont_from`.

```sh
docker rm cont_from
```

- **Expected Result**: The container `cont_from` will disappear from the list.

### 2. Removing Multiple Stopped Containers by ID

- **Command**: `docker rm <container_id1> <container_id2> ...`
- **Example**: Remove containers with specific IDs.

```sh
docker rm <container_id1> <container_id2>
```

- **Expected Result**: The specified containers will be removed and will no longer appear in the list.

### 3. Removing a Running Container with Force

Attempting to remove a running container without force will prompt an error:

```sh
docker rm <running_container_name>
```

- **Expected Error**: Docker will ask to use the `--force` flag.

To force remove a running container:

- **Command**: `docker rm -f <running_container_name>`
- **Example**: Remove a running container named `running_container`.

```sh
docker rm -f running_container
```

- **Expected Result**: The running container will be forcefully removed and will disappear from the list.

### 4. Killing a Running Container

If you prefer to kill a container before removing it:

- **Command**: `docker container kill <container_name>`
- **Example**: Kill a running container named `running_container`.

```sh
docker container kill running_container
```

- **Expected Result**: The container will be stopped.

### 5. Removing Stopped Containers with Prune

To remove all stopped containers:

- **Command**: `docker container prune`
- **Explanation**: This command removes all stopped containers, freeing up resources.

```sh
docker container prune
```

- **Expected Result**: All stopped containers will be removed, and Docker will report the amount of space freed up.

## Summary

We have covered various methods to remove Docker containers:

1. **Removing a Single Stopped Container by Name**:
    ```sh
    docker rm cont_from
    ```

2. **Removing Multiple Stopped Containers by ID**:
    ```sh
    docker rm <container_id1> <container_id2>
    ```

3. **Force Removing a Running Container**:
    ```sh
    docker rm -f running_container
    ```

4. **Killing a Running Container**:
    ```sh
    docker container kill running_container
    ```

5. **Removing All Stopped Containers with Prune**:
    ```sh
    docker container prune
    ```

By using these commands, you can efficiently manage and clean up your Docker containers, ensuring your environment remains organized and resource-efficient. In the next module, we will delve deeper into Docker networking.