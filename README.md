exporters
---

custmized dynamic exporter
   
usage:
```shell
curl http://x.x.x.x:6000/scrape?target=1.1.1.1:6379&service=redis
curl http://x.x.x.x:6000/scrape?target=1.1.1.1:9200&service=elasticsearch
curl http://x.x.x.x:6000/scrape?target=1.1.1.1:3306&service=mysql
```

currently only support:
* redis
* mysql
* elasticsearch
* logstash
* rocketmq
* kafka
  
   

auth info must provided within config file

---

working with prometheus requires `relabel_config` to change `[__address__]` label, a valid scrape config would like this:
```conf
  - job_name: 'consul-tencent'
    scrape_interval: 15s
    scrape_timeout: 10s
    metrics_path: '/scrape'
    honor_labels: false
    consul_sd_configs:
      - server: '127.0.0.1:8500'
        datacenter: 'dc1'
        services:
          - 'Elasticsearch'
          - 'MySQL'
          - 'Redis'
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__meta_consul_service]
        target_label: __param_service
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 127.0.0.1:6000
```
