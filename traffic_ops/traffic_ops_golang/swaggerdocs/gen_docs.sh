#!/usr/bin/env bash

unset DEBUG
#export DEBUG=true
swagger generate spec -o ./swagger.json
echo "successfully generated the swagger.json file"
