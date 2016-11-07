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
TO_UT_PROJ=traffic_ops_unit_tests
TO_HOME=/vol1/jenkins/jobs/$TO_UT_PROJ
TRAFFIC_OPS_APP=$TO_HOME/workspace/traffic_ops/app
export INSTALL_DIR=/opt/traffic_ops
export LOCAL_DIR=$INSTALL_DIR/app/local/lib
export PERL5LIB=$TRAFFIC_OPS_APP/lib:$LOCAL_DIR:$LOCAL_DIR/perl5/
export PATH=$PATH:$INSTALL_DIR/install/bin
#For Mojo
export MOJO_MODE=test

rm -fv $TRAFFIC_OPS_APP/*.tap
/bin/mkdir -p $TRAFFIC_OPS_APP/log
cd $TRAFFIC_OPS_APP
db/admin.pl --env=$MOJO_MODE setup

cd $TRAFFIC_OPS_APP

#Attempt to speed up the carton installs
#mkdir -p /tmp/$TO_UT_PROJ/local
#mkdir -p local/lib/perl5
#cp -r /tmp/$TO_UT_PROJ/local/* local/lib/perl5
#/usr/local/bin/carton

# Setup for the PERL5LIB
#cp -r local/lib/perl5/* lib

$INSTALL_DIR/app/local/bin/prove -vrp --formatter=TAP::Formatter::Jenkins t/

#cp -r local/lib/perl5 /tmp/$TO_UT_PROJ/local
