# speedtest-exporter
A prometheus exporter wrapping the [speedtest-cli](https://github.com/sivel/speedtest-cli). The speedtest-exporter exports the following metrics and is therfore compatible with any grafana dashboard created for [speedtest_exporter](https://github.com/nlamirault/speedtest_exporter)


## Build speedtest exporter

```bash
docker build -t speedtest-exporter . --build-arg ARCH=amd64 # use different architecture 
```

## Run speedtest exporter

```bash
docker run -p 9112:9112 --rm speedtest-exporter
```

Run a measurement at http://localhost:3000/metrics.

## Why 

I wanted to track my download and uplaod speed in my kubernetes homelab. I started to use the speedtest_exporter from nlamirault which was working very well for me - but I noticed that the measurements were quite off to browser based measurements. As I followed on the underlying dependencies I found that the speedtest implementation by github.com/zpeters/speedtest has been archived. Since I haven't found any prometheus exporter based on a maintained speedtest implementation I opted to wrap the unofficial speedtest-cli. 

## References

Many thanks to the work in the following respositories this speedtest prometheus exporter is based on.

* https://github.com/sivel/speedtest-cli
* https://github.com/nlamirault/speedtest_exporter
