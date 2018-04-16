#!/usr/bin/env bash

rm ./swagger.json
docker build -t tc-swaggerdocs -f Dockerfile-swagger-gen .
docker run --rm -it -v `(pwd)`:/output tc-swaggerdocs
