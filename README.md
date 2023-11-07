# Docker exporter

## Description

A small prometheus exporter used to monitor what cAdvisor does not. It currently supports :

- Docker images :
    - Name
    - Tag
    - Size
    - Number of containers using it
- Docker volumes :
    - Name
    - Mount point
    - Name of docker compose project
    - Number of containers using it
- Docker networks :
    - Name
    - Number of containers using it

## Usage

### Volumes

The exporter needs access to the docker socket to get the data. You can mount it inside the container :

```bash
docker run -v /var/run/docker.sock:/var/run/docker.sock cefadrom/docker-exporter:2
```

If you want to get the volumes size, you also need to mount the root filesystem inside the container :

```bash
docker run -v /var/run/docker.sock:/var/run/docker.sock -v /var/lib/docker/volumes:/var/lib/docker/volumes:ro cefadrom/docker-exporter:2
```

Note: for security reasons, it is recommended to mount data volumes read-only (using `:ro`).

### Docker compose

The recommended way is to use docker compose. Put the exporter config inside the same `docker-compose.yml` file as your
monitoring stack :

```yaml
services:
    prometheus:
      image: prom/prometheus:latest
        container_name: prometheus
        volumes:
            - ./prometheus.yml:/etc/prometheus/prometheus.yml
        ports:
            - 9090:9090

    docker-exporter:
      image: cefadrom/docker-exporter:2
        container_name: docker-exporter
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock # To access docker API
            - /var/lib/docker/volumes:/var/lib/docker/volumes:ro # If you want to get volumes size
        ports:
            - 9100:9100
```

Finally, configure your `prometheus.yml` file to scrape the exporter. Note that a high scrape interval is recommended as
the volume size calculation may be quite ressource intensive, and the data does not change that often.

```yaml
scrape_configs:
    -   job_name: 'docker-exporter'
        scrape_interval: 10m
        static_configs:
            -   targets: [ 'docker-exporter:9100' ]
```

## Metrics

Here is a sample of the metrics exposed by the exporter :

```text
# HELP docker_image_size_bytes Size of docker images
# TYPE docker_image_size_bytes gauge
docker_image_size_bytes{id="sha256:efadb309f7ffb875ad073162932b72d0bba5f8da71de54c36055fb843ef31be5",tag="cefadrom/docker-exporter:latest",containers="1"} 217327756
docker_image_size_bytes{id="sha256:3a8d46c63628f568039abf84f17c558cdaf4912d0c54a904e105f93f37f25775",tag="redis:alpine",containers="2"} 37773458

# HELP docker_volume_size_bytes Size of docker volumes
# TYPE docker_volume_size_bytes gauge
docker_volume_size_bytes{name="portainer_portainer_data",mountpoint="/var/lib/docker/volumes/portainer_portainer_data/_data",compose_project="portainer",containers="1"} 573341
docker_volume_size_bytes{name="dockprom_grafana_data",mountpoint="/var/lib/docker/volumes/dockprom_grafana_data/_data",compose_project="dockprom",containers="1"} 1915362
docker_volume_size_bytes{name="dockprom_prometheus_data",mountpoint="/var/lib/docker/volumes/dockprom_prometheus_data/_data",compose_project="dockprom",containers="1"} 1918972991

# HELP docker_network_container_count Number of containers per network
# TYPE docker_network_container_count gauge
docker_network_container_count{id="39b29ca27bc870e8259451db4c13b76ca88a0d211d3fe102eb5689938f094355",name="host"} 0
docker_network_container_count{id="ea78f0a6ef4949b684c155c5b11addb9cd79b7e2c0f7d43791c699373c5ca921",name="proxy"} 3
docker_network_container_count{id="d199cf3d3d45932efc1140e12922b07e8b31ba2715c1f686d44eff19ad7b69c7",name="bridge"} 0
docker_network_container_count{id="e9dcbdd6749b380fcc8789abbe7fe4eeb985110287d6ad4ef71026640688a914",name="none"} 0

# HELP docker_last_scrape_duration_seconds Time it took to scrape docker metrics
# TYPE docker_last_scrape_duration_seconds gauge
docker_last_scrape_duration_seconds 2.755
```
