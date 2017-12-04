#!/bin/bash#
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

if [[ -e /opt/tomcat ]]
then
  echo "/opt/tomcat installed"
else
  echo "Installing tomcat"
  cd /tmp
  wget http://archive.apache.org/dist/tomcat/tomcat-6/v6.0.33/bin/apache-tomcat-6.0.33.tar.gz
  cd /opt
  tar xfz /tmp/apache-tomcat-6.0.33.tar.gz
  ln -s apache-tomcat-6.0.33 tomcat
  ls -l /opt/tomcat/webapps/*
fi      

rm -rf /opt/tomcat/webapps/*
rm -f /opt/tomcat/bin/*.bat
chmod +x /opt/tomcat/bin/*.sh
