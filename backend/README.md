# Installation d'InfluxDB via Docker

## Commande Docker

Exécutez la commande suivante dans votre terminal :

```bash
docker run -d \
  --name influxdb \
  --restart always \
  -p 8086:8086 \
  -e DOCKER_INFLUXDB_INIT_MODE=setup \
  -e DOCKER_INFLUXDB_INIT_USERNAME=admin \
  -e DOCKER_INFLUXDB_INIT_PASSWORD=admin1234 \
  -e DOCKER_INFLUXDB_INIT_ORG=ntc-org \
  -e DOCKER_INFLUXDB_INIT_BUCKET=ntc-bucket \
  -v $(pwd)/influxdb-data:/var/lib/influxdb2 \
  -v $(pwd)/config:/etc/influxdb2 \
  influxdb:latest

```
