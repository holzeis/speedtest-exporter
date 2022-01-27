# speedtest-exporter
A prometheus export wrapping the speedtest-cli


## Build speedtest

```bash
docker build -t speedtest . --build-arg ARCH=amd64 # use different architecture 
```

## References

https://github.com/sivel/speedtest-cli
https://github.com/nlamirault/speedtest_exporter
