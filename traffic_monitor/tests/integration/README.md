# Traffic Monitor Integration Test Framework

## Building 

From the `trafficcontrol/traffic_monitor` directory run `build_tests.sh`

## Running
From the `trafficcontrol/traffic_monitor` directory run:

`sudo docker-compose -p tmi --project-directory . -f tests/integration/docker-compose.yml run tmintegrationtest`
