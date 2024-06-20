# Docker Volume Management

In this demo, we will walk through the process of creating, listing, inspecting, and removing Docker volumes using the command line.

## Step 1: Creating Volumes

First, we create a volume named `vol-busybox` using the `docker volume create` command.

```
docker volume create vol-busybox
```

Once the command succeeds, the name of the volume confirms its creation. Next, we create another volume by running a container using the Ubuntu image and mounting a volume named `vol-ubuntu` to the container's `/tmp` directory.

```
docker run -d --name tender_noyce -v vol-ubuntu:/tmp ubuntu
```

At this point, we have created two volumes: `vol-busybox` and `vol-ubuntu`.

## Step 2: Listing Volumes

To see the volumes created, we use the `docker volume ls` command.

```
docker volume ls
```

This command lists all the volumes, including those created by Docker using the local volume driver.
## Step 3: Filtering Volumes

We can filter the volumes to show only those that are not currently mounted to any container by using the filter `dangling=true`.

```
docker volume ls -f dangling=true
```

This command lists the volumes `vol-busybox` and another volume created by Docker that are not in use.

## Step 4: Inspecting Volumes

To get detailed information about a volume, we use the `docker volume inspect` command.

```
docker volume inspect vol-busybox
```

This command provides details such as the creation timestamp, driver type, mount point, and scope of the volume.


## Step 5: Removing Volumes

We attempt to remove a volume using the `docker volume rm` command.

```
docker volume rm vol-ubuntu
```

If the volume is in use, Docker returns an error. To find the container using the volume, we list the containers.

```
docker ps -a
```

Identifying the container `tender_noyce`, we remove it.

```
docker container rm tender_noyce
```

After removing the container, we can successfully remove the volume.

```
docker volume rm vol-ubuntu
```

To verify, we list the volumes again to ensure `vol-ubuntu` has been removed.

## Summary

In this demo, we created, listed, inspected, and removed Docker volumes. We saw how volumes can be managed and how they interact with containers, emphasizing the importance of proper volume management to ensure data persistence and application performance.

By following these steps and understanding Docker's volume management commands, you can efficiently handle storage in your Docker environment, ensuring data integrity and security.