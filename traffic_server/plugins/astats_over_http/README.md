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

Compile:
  tsxs -c astats_over_http.c -o astats_over_http.so
Install:
  sudo tsxs -o astats_over_http.so -i

Add to the plugin.conf:

  astats_over_http.so path=${path}

start traffic server and visit http://[ip]:[port]/${path}

Rpm Builds

  Two spec files are provided.  astats_over_http.spec requires a tar ball of this directory 
  named astats_over_htt-.tar.gz is copied to the rpmbuild/SOURCES directory.  The second
  astats-git-build, checks out the source from the git repo and builds the rpm.
