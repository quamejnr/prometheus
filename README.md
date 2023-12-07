## What is this?
A simple GRPC server with prometheus tracking


### Start grafana server

```
docker run --publish 3000:3000 --network host grafana/grafana
```

### Start Prometheus server

```
docker run \
    -p 9090:9090 \
    -v ./prometheus.yml:/etc/prometheus/prometheus.yml \
    --network host prom/prometheus
```

### NB:

- Add Prometheus as a data source in grafana by user the hos `http://localhost:9090`
- QUerying basics - https://prometheus.io/docs/prometheus/latest/querying/basics/
