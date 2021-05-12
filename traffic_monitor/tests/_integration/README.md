# Traffic Monitor Integration Test Framework

## Building

From the `trafficcontrol` run `./pkg -v traffic_monitor_build` and copy the traffic_monitor rpm to the `traffic_monitor` directory.

From the `trafficcontrol/traffic_monitor` directory run `build_tests.sh`

## Running

From the `trafficcontrol/traffic_monitor` directory run:

`sudo docker-compose -p tmi --project-directory . -f tests/_integration/docker-compose.yml run tmintegrationtest`
