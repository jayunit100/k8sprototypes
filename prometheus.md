run prometheus
`ce4f5d199102        prom/prometheus     "/bin/prometheus --c…"   About an hour ago   Up About an hour    0.0.0.0:9090->9090/tcp   elegant_cray`


with the following `-v` scrape target
```
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: 'codelab-monitor'

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'
    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s
    static_configs:
      - targets: ['10.0.0.217:2381']
      - targets: ['10.0.0.251:2381']
      - targets: ['10.0.0.141:2381']
```

where etcd is started with

```
    - --listen-metrics-urls=http://127.0.0.1:2381
```

Then ssh port forward into your bastion host (which has access to the 10.... addresses)

```
ssh -L 8080:127.0.0.1:9090 ubuntu@34.221.173.93
```

And browse metrics for all etcd hosts on 

`localhost:8080`
