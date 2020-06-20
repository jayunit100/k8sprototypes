mkdir /var/datap ; chmod 775 /var/datap   <1>
docker run -t -i -p 9090:9090 -v /var/datap:/var/datap -v
/tmp/p.yml:/etc/prometheus/prometheus.yml prom/prometheus
--storage.tsdb.path=/var/datap --config.file=/etc/prometheus/prometheus.yml <2>


cat << EOF > /tmp/prom.yml
global:
  scrape_interval:     3s <3> 
  external_labels:
    monitor: 'codelab-monitor'

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: 'prometheus'
    static_configs:
      - targets: ['10.0.0.217:2381'] <4>
      - targets: ['10.0.0.251:2381']
EOF


# Tools like `dd` and `fallocate -l 1G test.img` can be used to do a quick smoke of how fast your disks are.
