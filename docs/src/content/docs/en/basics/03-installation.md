---
title: Installation
description: Installation explanations.
---
# Running in Docker

Syncra provides official Docker images via Docker Hub that can be used for deployment on any system running Docker.

:::info
If you only plan to use the build-in executors, `http` and `shell` you can use the Syncra Light edition that only includes a single binary as the plugins are build-in.
:::

## Launching Syncra as a new container

Here’s a quick one-liner to get you off the ground (please note, we recommend further configuration for production deployments below):

```
docker run -d -p 8080:8080 --name taskvault danluki/taskvault agent --server --bootstrap-expect=1 --node-name=node1
```

Navigate to http://localhost:8080/ui

This will launch a Syncra server on port 8080 by default. You can use `docker logs -f syncra` to follow the rest of the initialization progress. Once the Syncra startup completes you can access the app at localhost:8080

Since Docker containers have their own ports, and we just map them to the system ports as needed it’s easy to move Syncra onto a different system port if you wish. For example running Syncra on port 12345:

```
docker run -d -p 12345:8080 --name taskvault danluki/taskvault agent --server --bootstrap-expect=1 --node-name=node1
```

## Mounting a mapped file storage volume

Syncra uses the local filesystem for storing the embedded database to store its own application data and the Raft protocol log. The end result is that your Syncra data will be on disk inside your container and lost if you ever remove the container.

To persist your data outside, of the container and make it available for use between container launches we can mount a local path inside our container.

```
docker run -d -p 8080:8080 -v ~/syncra.data:/syncra.data --name taskvault danluki/taskvault agent --server --bootstrap-expect=1 --data-dir=/syncra.data
```

Now when you launch your container we are mounting that folder from our local filesystem into the container.