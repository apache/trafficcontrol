#!/usr/bin/env bash

#    Licensed to the Apache Software Foundation (ASF) under one
#    or more contributor license agreements.  See the NOTICE file
#    distributed with this work for additional information
#    regarding copyright ownership.  The ASF licenses this file
#    to you under the Apache License, Version 2.0 (the
#    "License"); you may not use this file except in compliance
#    with the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#    Unless required by applicable law or agreed to in writing,
#    software distributed under the License is distributed on an
#    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
#    KIND, either express or implied.  See the License for the
#    specific language governing permissions and limitations
#    under the License.

# Uncomment (or set in your environment) to enable debug output for the swagger generation
#export DEBUG=true

OUTPUT_DIR=swaggerspec
SWAGGER_SPEC_FILE=$OUTPUT_DIR/swagger.json
mkdir -p $OUTPUT_DIR
swagger generate spec -o $SWAGGER_SPEC_FILE
echo "successfully generated swagger output file: $SWAGGER_SPEC_FILE"
