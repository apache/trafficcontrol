# Traffic Monitor Integration Test Framework

## Running

From the `trafficcontrol/traffic_monitor` directory:

```
(cd tools/testto && go build)
(cd tools/testcaches && go build)
(cd tests/integration && go test -c -o traffic_monitor_integration_test)
sudo docker-compose -p tmi --project-directory . -f tests/integration/docker-compose.yml run tmintegrationtest
```
