#!/bin/bash

# Update system shared CA support
source /to-access.sh

# Wait on SSL certificate generation
until [ -f "$CERT_DONE_FILE" ] 
do
     echo "Waiting on Shared SSL certificate generation"
     sleep 3
done

# Source the CIAB-CA shared SSL environment
source "$CERT_ENV_FILE"

# Copy the CIAB-CA certificate to the traffic_router conf so it can be added to the trust store
cp $CERT_CA_CERT_FILE /usr/local/share/ca-certificates
update-ca-certificates

# Grep out the existing SSL and Socket listener config
cp -af /etc/riak/riak.conf /etc/riak/riak.conf.orig
grep -v -E '^(listener|#)' /etc/riak/riak.conf.orig  | uniq | sort > /etc/riak/riak.conf

# Update the riak listener config
echo "nodename = riak@0.0.0.0" >> /etc/riak.conf
echo "listener.protobuf.internal = 0.0.0.0:$TV_INT_PORT" >> /etc/riak/riak.conf
echo "listener.http.internal = 0.0.0.0:$TV_HTTP_PORT" >> /etc/riak/riak.conf
echo "listener.https.internal = 0.0.0.0:$TV_HTTPS_PORT" >> /etc/riak/riak.conf

# Update SSL/TLS Certificate Config
echo "ssl.certfile = $CERT_TRAFFICVAULT_CERT" >> /etc/riak/riak.conf
echo "ssl.keyfile = $CERT_TRAFFICVAULT_KEY" >> /etc/riak/riak.conf
echo "ssl.cacertfile = /etc/pki/tls/certs/ca-bundle.crt" >> /etc/riak/riak.conf

# Enable search with Apache Solr
echo "search = on" >>  /etc/riak/riak.conf
