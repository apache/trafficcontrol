#!/bin/bash
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
TO_HOME=../../..
TRAFFIC_OPS_API_TEST=$TO_HOME/.
TRAFFIC_OPS_APP=$TO_HOME/traffic_ops/app
export INSTALL_DIR=/opt/traffic_ops
export LOCAL_DIR=$INSTALL_DIR/app/local/lib
export PERL5LIB=$TRAFFIC_OPS_APP/lib:$LOCAL_DIR:$LOCAL_DIR/perl5/
export PATH=$PATH:$INSTALL_DIR/install/bin
export CONFIG_HOME=.
#For Mojo
export MOJO_MODE=test

/bin/mkdir -p $TRAFFIC_OPS_APP/log
cd $TRAFFIC_OPS_APP
echo "TRAFFIC_OPS_APP: $TRAFFIC_OPS_APP"
db/admin.pl --env=$MOJO_MODE reset

cd $TRAFFIC_OPS_API_TEST
go test -v -cfg=$CONFIG_HOME/traffic-ops-test.conf -run TestCDNs
