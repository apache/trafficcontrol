# Private VNC/proxy testclient for ApacheConf Demo

## To start/install the VNC/proxy container

```
alias vdc='docker-compose -f docker-compose.yml -f docker-compose.testclient.yml'
vdc build 
vdc kill && vdc rm -f 
docker volume prune -f
vdc up
```

## Install a VNC client from the apple store or with brew
- host/port `localhost:55900`
- Password for VNC session is 'demo' or whatever $VNC_PASSWD is set to in Dockerfile

## URL locations within the VNC/proxy container:
* Traffic Portal: https://trafficportal.infra.ciab.test
* Traffic Monitor: http://trafficmonitor.infra.ciab.test
* Demo1 Delivery Service: http://video.demo1.mycdn.ciab.test/index.html

## TODO:
* On both linux/osx platforms, allow connectivity with docker daemon within the container
* Generate a chopped up HLS movie
