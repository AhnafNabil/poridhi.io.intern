# Docker Compose: Working with Compose Files

## Introduction

In this document, we will shift back to working with commands and files together, focusing on Docker Compose. We'll understand YAML files, and define services and volumes in a Docker Compose file.

## Understanding YAML Files

Before diving into the `docker-compose.yaml` file, it's essential to understand a few basics about YAML (YAML Ain't Markup Language) files. YAML files are used to define configurations in a human-readable format. They have three basic data types:

1. **Scalars**: Strings and Numbers
2. **Sequences**: Arrays or Lists
3. **Mappings**: Hashes or Dictionaries, represented as key-value pairs

Nesting of objects in YAML is determined by indentation.

## The Docker Compose File

The `docker-compose.yaml` file defines multiple objects such as services, networks, and volumes. The default path for the compose file is always the present directory. 

### Defining the Version

First, we need to specify the version of Docker Compose we are using. In this case, we will use version 3.3.

```yaml
version: '3.3'
```

### Defining Services

Services are the parent object for the containers we will create. For a multi-container application, we define each container under the `services` section.

#### Creating the `db` Service

Let's create our first service, `db`, which stands for the database. We need to specify a few parameters as key-value pairs:

- **Image**: We will use MySQL version 5.7.
- **Container Name**: We will name the container `mysql_database`.
- **Volumes**: We define the volume name and mount path.
- **Restart Policy**: We set this to `always`.
- **Environment Variables**: We define necessary environment variables for MySQL.

```yaml
services:
  db:
    image: mysql:5.7
    container_name: mysql_database
    volumes:
      - db_data:/var/lib/mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: word@press
      MYSQL_DATABASE: wordpress
      MYSQL_USER: wordpress
      MYSQL_PASSWORD: abc@123
```

#### Creating the `wordpress` Service

Next, we create the `wordpress` service. This service depends on the `db` service, meaning the `db` container must be created first.

- **Image**: We will use the WordPress image.
- **Container Name**: We will name the container `wd_frontend`.
- **Volumes**: We define the volume name and mount path.
- **Ports**: We map port 8000 on the host to port 80 in the container.
- **Restart Policy**: We set this to `always`.
- **Environment Variables**: We define necessary environment variables for WordPress.

```yaml
  wordpress:
    depends_on:
      - db
    image: wordpress
    container_name: wd_frontend
    volumes:
      - wordpress_files:/var/www/html
    ports:
      - "8000:80"
    restart: always
    environment:
      WORDPRESS_DB_HOST: db:3306
      WORDPRESS_DB_USER: wordpress
      WORDPRESS_DB_PASSWORD: abc@123
```

### Defining Volumes

Finally, we define the volumes used in the services. These are declared outside the `services` field.

```yaml
volumes:
  wordpress_files: 
  db_data: 
```

## Summary

In summary, we have created a Docker Compose file that defines two services: `db` (MySQL) and `wordpress`. We have specified all necessary fields for each service, including container name, image, environment variables, and volume mount information. Additionally, we have defined the volumes used by these services.

In the next document, we will execute this Compose file and observe how the application works.