<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

Traffic Control in Docker
=========================

Traffic Control wasn't designed to run in containers or with the microservice ideology, so there's a bit of coercing to make Docker work. But feel free to try it if you're feeling adventurous.

IP addresses and hostnames
--------------------------
The CDN requires real hostnames. Ideally, internal CDN communication could take place via Docker networking, using Docker internal IPs, hostnames, and container names. But the CDN doesn't currently distinguish between internal and external communication addresses. Further, the version of Apache Traffic Server we're using doesn't use `/etc/hosts`, which is what Docker uses to share hostnames between containers on a Docker Network.

So, we typically allocate public hostnames, domains, and IPs for each container. See the examples below.

Note if you have multiple IP addresses on a single network interface, you will have to create IP aliases. Docker requires a single interface per IP. You can create an Ethernet alias in Linux with `ifconfig` with, for example, `sudo ifconfig myinterface0:1 192.0.2.0 netmask 255.255.255.0 up`, where `myinterface0` is the interface with a range of IP addresses, and `192.0.2.0` is the IP you want to create the alias for. You can then forward a port on a Docker container with, for example, `docker run --publish 192.0.2.0:443:443`.


Docker Networks
-----------------
[Docker Networks](https://docs.docker.com/engine/userguide/networking/dockernetworks/), as opposed to `--link` or other container communication mechanisms, is necessary for a number of reasons.

Most of the following is only necessary if you're trying to set up a CDN in Docker without external IPs or hostnames. But these things seem to be generally good practice.

#### Container names

- Containers MUST have names which are valid hostnames
- Containers MUST have the same name and hostname. That is, `docker run` `--name` and `--hostname` MUST match.

This is because `docker run` `--hostname` sets the container's `hostname` ($HOSTNAME), while `--name` sets the /etc/hosts entry for _other_ containers. If they are different, the name by which a container knows itself will be different from the name by which other containers access it. This will lead to confusion and pain.

#### Hostnames

Docker Networks automatically update every container on the network's `/etc/hosts` file with the hostname of a new container. Without this, for example, you would have to somehow manually update the Traffic Ops container when bringing up a new Traffic Vault container, with the new container's hostname–IP mapping. This is strictly necessary because Traffic Control requires hostnames and domains, not IPs, in certain places.

#### Domains

Docker networks also updates the `/etc/hosts` file of member containers with `hostname.networkname`, which should theoretically allow us to treat the Docker network name as the domain name within Traffic Control. For example, you will see the network name `cdnet` used in the example run command in the Dockerfiles, and then `cdnet` could also (theoretically) be passed as the domain name to various Traffic Control components.

Without Docker Networks, each container's hostname could be set as `hostname.domainname`, but that would cause other issues, such as `hostname` itself not existing as a host name.

Without Docker Networks, each Traffic Control component must also have to have its ports exposed on the host machine, or else use the deprecated Docker `--link` mechanism.

Run scripts
-----------

Each Dockerfile creates a generic image, which should be usable in any environment. If a Dockerfile includes CDN-specific values or data, that's a bug. Consequently, `docker build` should require few, if any, environment variables.

Each container has a `_run.sh` script which configures the container. While self-contained Dockerfiles would be ideal, in practice the Traffic Control components need dynamic configuration. If self-contained Dockerfiles were necessary, an `echo` command could be added to the Dockerfile to create the run script; but it would be ugly.

Each run script has 2 functions, `init` and `start`. The init is run only once, and configures the container's service. After the first run, only `start` is run.

Cache devices
-------------

The Apache Traffic Server containers look for devices at `/dev/ram0` and `/dev/ram1`. You can pass devices on the host to containers with `docker run` `--device`. If no devices are passed, the run script will try to create 1GB files to use as caches.

You can create block RAM devices on Linux with, for example, `sudo modprobe brd rd_size=1048576 rd_nr=8`. This will create 8 RAM disks at `/dev/ram0` thru `/dev/ram7` of 1GB each. These devices can then be passed to the Apache Traffic Server containers, as in the example below.

Privileged containers
---------------------

The Traffic Server containers must be run with `--privileged` and `--cap-add NET_BIND_SERVICE` in order for Apache Traffic Server to bind on port 80. If you don't want to run privileged containers, you can try modifying the containers to run on another port (you will also need to update their ports in the Traffic Ops server entry, update the Traffic Ops `server_ports` parameters, re-generate the CRConfig in Traffic Ops, and re-run the `ort` script on the caches).

Example
--------

Suppose you have the following IPs on the host, each on their own network interface (or alias), with the following hostnames

IP          | Hostname
------------|---------
192.0.2.100 | c23-to-db
192.0.2.101 | c23-to-01
192.0.2.102 | c23-tv-01
192.0.2.103 | c23-hotair-01
192.0.2.104 | c23-atsec-01
192.0.2.105 | c23-atsec-01
192.0.2.106 | c23-atsmid-01
192.0.2.107 | c23-atsmid-02
192.0.2.108 | c23-tm-01
192.0.2.109 | c23-tm-02
192.0.2.110 | c23-tr-01
192.0.2.111 | c23-ts-01

And the NS record `c23` registered with `example.net`.

You will also need an Origin server. This just needs to be a simple HTTP server, serving some content somewhere. For this example, we will assume there is an origin server with the IP `192.0.2.103` and hostname `c23-hotair-01` serving `http://c23-hotair-01.example.net/test.ism/manifest`.

Finally, for this example, we will assume the host has cache drives at `/dev/ram0` thru `/dev/ram7`. If you don't want to pass drives, simply omit the `--device` flags from the examples, and the containers will try to create disk files to use as cache instead.

In the example commands, `-it` runs the container interactively, so you can watch the initialization process, and (hopefully) see if anything goes wrong. Once the initialization finishes, you can detach from the container with `ctrl+p ctrl+q`. Alternatively, remove the `-it` flags and add `-d` to run the containers detached. However, certain containers need to finish initializing before the next container can be started, and without running interactively it can be difficult to tell.

With the prior assumptions, the following commands will set up a CDN:

```bash
docker network create cdnet

sudo docker run -it --publish 192.0.2.100:3306:3306 --name c23-to-db --hostname c23-to-db --net cdnet --env MYSQL_ROOT_PASS=secretrootpass --env IP=192.0.2.100 mysql:5.5

sudo docker run -it --publish 192.0.2.101:443:443 --name c23-to-01 --hostname c23-to-01 --net cdnet --env MYSQL_IP=c23-to-db.example.net --env MYSQL_PORT=3306 --env MYSQL_ROOT_PASS=secretrootpass --env MYSQL_TRAFFIC_OPS_PASS=supersecretpassword --env ADMIN_USER=superroot --env ADMIN_PASS=supersecreterpassward --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --env TRAFFIC_VAULT_PASS=marginallylesssecret --env IP=192.0.2.101 --env DOMAIN=c23.example.net traffic_ops:1.4

sudo docker run -it --publish 192.0.2.102:8088:8088 --name c23-tv-01 --hostname c23-tv-01 --net cdnet --env ADMIN_PASS=riakadminsecret --env USER_PASS=marginallylesssecret --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --env TRAFFIC_OPS_URI=https://c23-to-01.example.net --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=example.net --env IP=192.0.2.102 --env GATEWAY=192.0.2.161 traffic_vault:1.4

sudo docker run -it --publish 192.0.2.104:80:80 --name c23-atsec-01 --hostname c23-atsec-01 --net cdnet --privileged --cap-add NET_BIND_SERVICE --device /dev/ram0:/dev/ram0 --device /dev/ram1:/dev/ram1 --env TRAFFIC_OPS_URI=https://c23-to-01.example.net --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=example.net --env IP=192.0.2.104 --env GATEWAY=192.0.2.161 traffic_server_edge:1.4

sudo docker run -it --publish 192.0.2.105:80:80 --name c23-atsec-02 --hostname c23-atsec-02 --net cdnet --privileged --cap-add NET_BIND_SERVICE --device /dev/ram2:/dev/ram0 --device /dev/ram3:/dev/ram1 --env TRAFFIC_OPS_URI=https://c23-to-01.example.net --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=example.net --env IP=192.0.2.105 --env GATEWAY=192.0.2.161 traffic_server_edge:1.4

sudo docker run -it --publish 192.0.2.106:80:80 --name c23-atsmid-01 --hostname c23-atsmid-01 --net cdnet --privileged --cap-add NET_BIND_SERVICE --device /dev/ram4:/dev/ram0 --device /dev/ram5:/dev/ram1 --env TRAFFIC_OPS_URI=https://c23-to-01.example.net --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=example.net --env IP=192.0.2.106 --env GATEWAY=192.0.2.161 traffic_server_mid:1.4

sudo docker run -it --publish 192.0.2.107:80:80 --name c23-atsmid-02 --hostname c23-atsmid-02 --net cdnet --privileged --cap-add NET_BIND_SERVICE --device /dev/ram6:/dev/ram0 --device /dev/ram7:/dev/ram1 --env TRAFFIC_OPS_URI=https://c23-to-01.example.net --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=example.net --env IP=192.0.2.107 --env GATEWAY=192.0.2.161 traffic_server_mid:1.4

sudo docker run -it --publish 192.0.2.108:80:80 --name c23-tm-01 --hostname c23-tm-01 --net cdnet --env TRAFFIC_OPS_URI=https://c23-to-01.example.net --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=example.net --env IP=192.0.2.108 --env GATEWAY=192.0.2.161 traffic_monitor:1.4

sudo docker run -it --publish 192.0.2.109:80:80 --name c23-tm-02 --hostname c23-tm-02 --net cdnet --env TRAFFIC_OPS_URI=https://c23-to-01.example.net --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env DOMAIN=example.net --env IP=192.0.2.109 --env GATEWAY=192.0.2.161 traffic_monitor:1.4

# The following commands will add the Origin server to Traffic Ops.
# If your Origin is built in Docker, you could add these commands to the container run script.
TMP_IP='192.0.2.103'
TMP_DOMAIN='example.net'
TMP_GATEWAY='192.0.2.161'
TMP_TRAFFIC_OPS_USER='superroot'
TMP_TRAFFIC_OPS_PASS='supersecreterpassward'
TMP_TRAFFIC_OPS_URI='https://c23-to-01.example.net'
TMP_HOSTNAME='c23-hotair-01'
TMP_TO_COOKIE="$(curl -v -s -k -X POST --data '{ "u":"'"$TMP_TRAFFIC_OPS_USER"'", "p":"'"$TMP_TRAFFIC_OPS_PASS"'" }' $TMP_TRAFFIC_OPS_URI/api/1.2/user/login 2>&1 | grep 'Set-Cookie' | sed -e 's/.*mojolicious=\(.*\); expires.*/\1/')"
TMP_CACHEGROUP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TMP_TRAFFIC_OPS_URI/api/1.2/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="mid-east"]; print match[0]')"
TMP_SERVER_TYPE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TMP_TRAFFIC_OPS_URI/api/1.2/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="ORG"]; print match[0]')"
TMP_SERVER_PROFILE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TMP_TRAFFIC_OPS_URI/api/1.2/profiles.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="ORG1_CDN1"]; print match[0]')"
TMP_PHYS_LOCATION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TMP_TRAFFIC_OPS_URI/api/1.2/phys_locations.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="plocation-nyc-1"]; print match[0]')"
TMP_CDN_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TMP_TRAFFIC_OPS_URI/api/1.2/cdns.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="cdn"]; print match[0]')"
curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "host_name=$TMP_HOSTNAME" --data-urlencode "domain_name=$TMP_DOMAIN" --data-urlencode "interface_name=eth0" --data-urlencode "ip_address=$TMP_IP" --data-urlencode "ip_netmask=255.255.0.0" --data-urlencode "ip_gateway=$TMP_GATEWAY" --data-urlencode "interface_mtu=9000" --data-urlencode "cdn=$TMP_CDN_ID" --data-urlencode "cachegroup=$TMP_CACHEGROUP_ID" --data-urlencode "phys_location=$TMP_PHYS_LOCATION_ID" --data-urlencode "type=$TMP_SERVER_TYPE_ID" --data-urlencode "profile=$TMP_SERVER_PROFILE_ID" --data-urlencode "tcp_port=8088" $TMP_TRAFFIC_OPS_URI/server/create
TMP_SERVER_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TMP_TRAFFIC_OPS_URI/api/1.2/servers.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["hostName"]=="'"$TMP_HOSTNAME"'"]; print match[0]')"
curl -v -k -H "Content-Type: application/x-www-form-urlencoded" -H "Cookie: mojolicious=$TMP_TO_COOKIE" -X POST --data-urlencode "id=$TMP_SERVER_ID" --data-urlencode "status=ONLINE" $TMP_TRAFFIC_OPS_URI/server/updatestatus

sudo docker run -it --publish 192.0.2.110:80:80 --publish 192.0.2.110:3333:3333 --publish 192.0.2.110:53:53 --publish 192.0.2.110:53:53/udp --name c23-tr-01 --hostname c23-tr-01 --net cdnet --env TRAFFIC_OPS_URI=https://c23-to-01.example.net --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env TRAFFIC_MONITORS="c23-tm-01.example.net:80;c23-tm-02.example.net:80" --env ORIGIN_URI="http://c23-hotair-01.example.net" --env DOMAIN=example.net --env IP=192.0.2.110 --env GATEWAY=192.0.2.161 traffic_router:1.4

sudo docker run -it --publish 192.0.2.111:8083:8083 --publish 192.0.2.111:8086:8086 --name c23-ts-01 --hostname c23-ts-01 --net cdnet --env TRAFFIC_OPS_URI=https://c23-to-01.example.net --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --env DOMAIN=example.net --env IP=192.0.2.111 --env GATEWAY=192.0.2.161 traffic_stats:1.4

# Edge and Mid caches must be created first to add themselves as servers to Traffic Ops.
# Once Traffic Ops has all the servers and the CRConfig is generated (by the Traffic Router container), we must re-run the ort script on them.
sudo docker exec -it c23-atsec-01 /opt/ort/traffic_ops_ort.pl badass WARN https://c23-to-01.example.net superroot:supersecreterpassward
sudo docker exec -it c23-atsec-02 /opt/ort/traffic_ops_ort.pl badass WARN https://c23-to-01.example.net superroot:supersecreterpassward
sudo docker exec -it c23-atsmid-01 /opt/ort/traffic_ops_ort.pl badass WARN https://c23-to-01.example.net superroot:supersecreterpassward
sudo docker exec -it c23-atsmid-02 /opt/ort/traffic_ops_ort.pl badass WARN https://c23-to-01.example.net superroot:supersecreterpassward
```

If everything was successful, you can test the CDN with:

```bash
# Test the Traffic Router DNS. Should return the DNS entries for the delivery services
dig @c23-tr-01.example.net edge.ds1.c23.example.net

# Test the DNS delivery service. Should return a 200
curl -vs -o /dev/null -H "Host: edge.ds1.c23.example.net" "http://192.0.2.104/test.ism/manifest"

# Test the HTTP delivery service. Should return a 302 and a 200.
curl -Lvs -o /dev/null "http://tr.ds2.c23.example.net/test.ism/manifest"
```


Traffic Portal
-------------

Traffic Portal is not required by the CDN to function, but provides a web interface to manage Traffic Ops.

You can also create a standalone Traffic Portal Docker container and point it at a Traffic Ops not in Docker.

You will need a Traffic Portal RPM, and a running Traffic Ops instance to point it to.

The following commands will build a Docker image and container for Traffic Portal:


```
docker build --no-cache --rm --tag traffic_portal:3.0.x --build-arg=RPM=traffic_portal.rpm .

docker run --name tp --hostname tp --net cdnet --publish 40443:443 --env TO_SERVER=my-traffic-ops-fqdn --env TO_PORT=443 --env DOMAIN=cdnet --detach -- traffic_portal:3.0.x
```
