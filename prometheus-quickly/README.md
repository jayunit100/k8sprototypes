Prometheus instructions get more and more complex every day, but if you just want to monitor ETCD, this is the easiest way, no load balancers or operators or anything else required.

# Monitoring ETCD comes down to one metric: FSync

Look at its value, if it flat lines, it means that the originally fast write speed you had is slowing down, and etcd is 
not able to keep up.

![Image description](graph.png)


# Requirements

- 3 node etcd cluster
- ability to edit the kubelet manifests/etcd.yaml (i.e. if starting in CAPI this is the norm)
- SSH into a bastion node that can access the 3 nodes via IP address
- Ability to run docker on the bastion

Expectations: 

# note that the slowest write is 1/4 of a second
![Image description](prometheus.png)

# comparitively, in a slow cluster (without SSDs, or hypervirtualized, for example) you may very early on see something like this:

```
# TYPE etcd_disk_wal_fsync_duration_seconds histogram   
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.001"} 0
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.002"} 0
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.004"} 202
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.008"} 1601
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.016"} 2173
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.032"} 2552                                                                 
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.064"} 2635
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.128"} 2658
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.256"} 2669
etcd_disk_wal_fsync_duration_seconds_bucket{le="0.512"} 2674 <-- 5 writes took a half second1
etcd_disk_wal_fsync_duration_seconds_bucket{le="1.024"} 2676
etcd_disk_wal_fsync_duration_seconds_bucket{le="2.048"} 2676
etcd_disk_wal_fsync_duration_seconds_bucket{le="4.096"} 2676
etcd_disk_wal_fsync_duration_seconds_bucket{le="8.192"} 2676
etcd_disk_wal_fsync_duration_seconds_bucket{le="+Inf"} 2676
etcd_disk_wal_fsync_duration_seconds_sum 33.182538455000035
etcd_disk_wal_fsync_duration_seconds_count 2676      
```
If you start out a cluster this way, then there is a chance over time that performance will degrade even more, rapidly approaching > 1 second writes. 
# run prometheus in docker
```
docker run -p 9090:9090 -v prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

with the following `-v` scrape target (prometheus.yml)... 


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


In general you should see most prometheus writes happening within .001 s or less, like so:
# start etcd like this so you can scrape w/o credentials 

Each individual etcd node needs to get modified to listen on an insecure port like this
```
    - --listen-metrics-urls=http://10.0.0.217:2381
```

Then ssh port forward into your bastion host (which has access to the 10.... addresses)

```
ssh -L 8080:127.0.0.1:9090 ubuntu@34.221.173.93
```

And browse metrics for all etcd hosts on 

`localhost:8080`
