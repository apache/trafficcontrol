#!/bin/bash

set -o posix
set -o errexit
set -e



die() {
  local MSG=$1
  local RETURN_CODE=$2

  if [ -z $RETURN_CODE ]; then
    RETURN_CODE="1"
  fi

  printf "\e[31m ERROR: %s\e[0m\n" "$MSG" >&2

  return $RETURN_CODE
}
upgrade(){
trsvr=$1

if [ -z "$trsvr" ]; then
  die "Invalid traffic server"
fi

ssh -q -o StrictHostKeyChecking=no $trsvr rm -f *.rpm

ssh -q -o StrictHostKeyChecking=no $trsvr sudo systemctl stop traffic_router
ssh -q -o StrictHostKeyChecking=no $trsvr mkdir last 
ssh -q -o StrictHostKeyChecking=no $trsvr sudo cp -f /opt/traffic_router/conf/*.properties last 
ssh -q -o StrictHostKeyChecking=no $trsvr sudo yum remove -y tomcat
ssh -q -o StrictHostKeyChecking=no $trsvr sudo rm -rf /opt/traffic_router
ssh -q -o StrictHostKeyChecking=no $trsvr sudo yum --enablerepo artifacts-nightly install -y traffic_router
ssh -q -o StrictHostKeyChecking=no $trsvr sudo cp -f last/*.properties /opt/traffic_router/conf
ssh -q -o StrictHostKeyChecking=no $trsvr sudo systemctl daemon-reload
ssh -q -o StrictHostKeyChecking=no $trsvr sudo systemctl start traffic_router
echo "Deployed Traffic Router On Server: " $trsvr
}


upgrade ict11-tr-01.cdnlab.comcast.net
upgrade ict11-tr-02.cdnlab.comcast.net
upgrade ict11-tr-03.cdnlab.comcast.net
upgrade ict11-tr-04.cdnlab.comcast.net
