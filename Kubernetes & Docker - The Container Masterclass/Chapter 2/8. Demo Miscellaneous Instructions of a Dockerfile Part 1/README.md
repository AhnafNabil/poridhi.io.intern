# Docker Multi-Stage Build Demonstration


## Overview
In this demonstration, we will explore the concept of using multiple Dockerfiles within a single directory. This approach, while sometimes questioned, is perfectly valid and useful in certain scenarios. We'll clarify why having multiple Dockerfiles isn't a problem and then walk through a practical example using two Dockerfiles: a parent and a child Dockerfile.

## Dockerfiles
We start in a directory, which contains the following files:

- parent-dockerfile
- child-dockerfile

## Clarification on Multiple Dockerfiles

### Common Concerns
- Why have multiple Dockerfiles in one directory?
- Isn't it bad practice?
- Wouldn't Docker get confused?

### Explanation
- **Multiple Dockerfiles**: It is possible to have more than one Dockerfile in a repository or folder, but they must have different names.
- **Naming**: Operating systems do not allow multiple files with the same name in a single directory. Naming a Dockerfile as Dockerfile has the sole purpose of making image building commands shorter, especially with Docker Hub auto builder.
- **Docker Behavior**: Docker does not get confused by multiple Dockerfiles. It simply builds the file specified in the build command.

## Detailed Walkthrough

### Parent Dockerfile
Let's start by examining the parent Dockerfile.

```dockerfile
FROM ubuntu:16.04

ONBUILD RUN echo "Greetings from your parent image!" > /tmp/hello.txt

CMD ["bash"]
```
- Base Image: ubuntu:16.04

- ONBUILD Instruction: This command is triggered when the image built from this Dockerfile is used as a base for another image. It creates a file greetings.txt in the /tmp directory with the content "Greetings from your parent image!".

- CMD Instruction: Sets the default command to bash.


### Child Dockerfile
Next, let's look at the child Dockerfile.

```dockerfile
FROM papa-ubuntu:latest

CMD ["bash"]
```
- Base Image: papa-ubuntu:latest, which will be the image built from the parent Dockerfile.
- CMD Instruction: Sets the default command to bash.


## Building and Running the Docker Images
### Building the Parent Image
We build the parent image using the following command:

```sh
docker build -f parent-dockerfile -t papa-ubuntu:latest .
```

### Building the Child Image
Next, we build the child image:

```sh
docker build -f child-dockerfile -t baby-ubuntu:latest .
```

During the build process, the `ONBUILD` instruction in the parent Dockerfile is triggered, ensuring that the `greetings.txt` file is created in the child image.

## Verifying the Images
To verify that both images have been created, we list the Docker images:

```sh
docker images
```

### Running the Child Image
We run a container from the baby-ubuntu image and name it baby-container:

```sh
docker run --name baby-container -it baby-ubuntu:latest
```

Once inside the container, we navigate to the `/tmp` directory and check for the`greetings.txt` file:

```sh
cd /tmp
cat greetings.txt
```

The output should display "Greetings from your parent image!", confirming the success of the `ONBUILD` instruction.

## Conclusion
Through this demonstration, we have shown that using multiple Dockerfiles in a single directory is feasible and can be beneficial. The `ONBUILD` instruction allows for the creation of base images that automatically execute specific commands when used as the foundation for other images. This setup helps in maintaining consistency and reducing repetitive tasks in Dockerfile configurations.