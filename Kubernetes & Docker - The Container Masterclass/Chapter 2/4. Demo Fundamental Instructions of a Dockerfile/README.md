# Creating Our First Dockerfile: Step-by-Step Guide

In this guide, we will create our first Dockerfile and explore its fundamental instructions. The Dockerfile is an essential component for defining how a Docker image is built. Follow these steps to get started:

## Step 1: Verify Current Directory

First, verify your current working directory. If we are in a user directory, such as `/home/nabil`.

```sh
pwd
```

## Step 2: Navigate to the `CC_Docker` Directory

Use the `ls` command to list the contents and ensure the `CC_Docker` directory is present. Then, navigate to it otherwise make a directory named `CC_Docker`.

```sh
ls
cd CC_Docker
```

## Step 3: Create an Empty Dockerfile

Create an empty Dockerfile using the `touch` command.

```sh
touch Dockerfile
```

## Step 4: Open and Edit the Dockerfile

Open the Dockerfile using a text editor such as `nano`. You can use any text editor you prefer.

```sh
nano Dockerfile
```

## Step 5: Write the Dockerfile

Start writing the Dockerfile with fundamental instructions:

1. **ARG Instruction:**
   The `ARG` instruction defines a variable that users can pass at build-time to the builder with the `docker build` command. This helps keep parameters like versions under control.

    ```Dockerfile
    ARG CODE_VERSION=16.04
    ```

2. **FROM Instruction:**
   The `FROM` instruction sets the base image for subsequent instructions. It is essential in every Dockerfile and typically follows the `ARG` instruction if one is used.

    ```Dockerfile
    FROM ubuntu:${CODE_VERSION}
    ```

3. **RUN and CMD Instructions:**
   We include basic `RUN` and `CMD` instructions to add functionality, although their detailed exploration will come later.

    ```Dockerfile
    RUN apt-get update -y
    CMD ["bash"]
    ```

## Example Dockerfile

Here is the complete Dockerfile:

```Dockerfile
ARG CODE_VERSION=16.04
FROM ubuntu:${CODE_VERSION}
RUN apt-get update -y
CMD ["bash"]
```

## Step 6: Save and Exit

Save the file and exit the text editor. In `nano`, we can do this by pressing `CTRL + O` to save and `CTRL + X` to exit.

## Step 7: Build the Docker Image

Use the `docker build` command to build the Dockerfile into an image. The `-t` option tags the image with a name for easy reference.

```sh
docker build -t img_from .
```

## Step 8: Verify the Image

Verify that the image has been created using the `docker images` command. This command lists all Docker images on the system.

```sh
docker images
```

We should see the newly created `img_from` image listed among other Docker images.

## Summary

- **ARG Instruction:** Defines build-time variables.
- **FROM Instruction:** Sets the base image.
- **RUN and CMD Instructions:** Add functionality to the image.

By following these steps, we have created a simple Dockerfile, built a Docker image, and verified its creation. In the next sections, we will delve deeper into the `RUN` and `CMD` instructions and explore their uses in Dockerfiles.