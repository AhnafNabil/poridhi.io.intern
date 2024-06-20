# Dockerfile Section: Final Lecture

## Goals
In this document, we will:
1. Learn how to add health checks to containers using Dockerfiles.
2. Use the stopsignal instruction.
3. Create a container for a simple Flask application.

## Our Files
We are working in a directory, which has three files:
- `app.py`
- `requirements.txt`
- `Dockerfile`


Let's look at each file, starting with `app.py`.

## Understanding app.py
This file contains a basic Flask application. Flask is a tool that helps Python programs work with web servers.

### What's in app.py?
```python
from flask import Flask

app = Flask(__name__)

@app.route('/')
def CMC():
    return 'Welcome to the Container Master Class by Cerulean Canvas'

if __name__ == '__main__':
    app.run(host='0.0.0.0')
```
- **Importing Flask**: We bring in the Flask class from the Flask library.
- **Creating an App**: We create an app using `Flask(__name__)`.
- **Setting a Route**: The `@app.route('/')` line tells our app to respond to requests to the main URL (`/`).
- **Defining a Function**: The `CMC` function sends back the message 'Welcome to the Container Master Class by Cerulean Canvas'.
- **Running the App**: If this file is run directly, it starts the Flask web server.

## Understanding requirements.txt
This file lists the libraries our app needs. It has one line:
```
Flask==0.12.2
```
This tells us we need Flask version 0.12.2.

## Understanding Dockerfile
The Dockerfile describes how to build our Docker image.

### What's in the Dockerfile?
```dockerfile
FROM ubuntu:16.04

RUN apt-get update -y && \
    apt-get install -y python-pip python-dev curl

COPY . /app

WORKDIR /app

RUN pip install -r requirements.txt

HEALTHCHECK --interval=10s --timeout=30s CMD curl --fail http://localhost:5000/ || exit 1

CMD ["python", "app.py"]

STOPSIGNAL SIGKILL
```
- **Base Image**: We start with Ubuntu 16.04.
- **Install Software**: We update and install Python, pip, and curl.
- **Copy Files**: We copy our files to the `/app` directory in the image.
- **Set Work Directory**: We set `/app` as our working directory.
- **Install Dependencies**: We install the libraries listed in `requirements.txt`.
- **Health Check**: 
  - **interval**: Check the container every 10 seconds.
  - **timeout**: Give the check 30 seconds to complete.
  - **CMD**: Run a command to check if the app is responding at `http://localhost:5000`. If it fails, the container is marked unhealthy.
- **Run the App**: We use `CMD` to run `app.py` with Python.
- **Stop Signal**: We set the stop signal to `SIGKILL` to force the container to stop immediately when we shut it down.

### Health Check Details
Health checks help us know if the container is working correctly. The `HEALTHCHECK` instruction runs a command to test this. If the command fails, Docker marks the container as unhealthy.

### Stop Signal Details
The `STOPSIGNAL` instruction tells Docker how to stop the container. By default, Docker uses `SIGTERM`, which stops the container gently. We change it to `SIGKILL` to stop it immediately.

## Building and Running the Docker Image
1. **Build the Image**:
   ```bash
   docker build -t flask-app .
   ```
   This command creates our Docker image.

2. **Run the Container**:
   ```bash
   docker run -d --name flask flask-app
   ```
   This command runs a container named `flask` from our image.

3. **Check Containers**:
   ```bash
   docker ps
   ```
   This shows our running containers. The `flask` container should be listed as running and healthy.

4. **Verify Health Check**:
   ```bash
   curl http://localhost:5000
   ```
   This command tests if our Flask app is running. It should return 'Welcome to the Container Master Class by Cerulean Canvas'.

5. **Stop the Container**:
   ```bash
   docker stop flask
   ```

6. **Check Stop Signal**:
   ```bash
   docker ps -a
   ```
   This shows all containers, including stopped ones. The `flask` container should show it was stopped with an error code 137, indicating it was stopped by `SIGKILL`.

## Conclusion
We have:
1. Added health checks to our Dockerfile.
2. Used the `STOPSIGNAL` instruction.
3. Created and ran a Flask application in a container.
