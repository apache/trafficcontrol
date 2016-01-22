
#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
chkconfig --add tomcat
chkconfig tomcat on

if [ ! -e /opt/tomcat/lib/traffic_router_connector.jar ]; then
	#echo "Symlinking /opt/traffic_router/lib/traffic_router_connector.jar to /opt/tomcat/lib/traffic_router_connector.jar"
	/bin/ln -s /opt/traffic_router/lib/traffic_router_connector.jar /opt/tomcat/lib/traffic_router_connector.jar
fi

for app in core api; do
	if [ ! -d /opt/traffic_router/webapps/${app} ]; then
		/bin/mkdir /opt/traffic_router/webapps/${app}
	fi

	if [ ! -e /opt/traffic_router/webapps/${app}/ROOT.war ]; then
		#echo "Symlinking /opt/traffic_router/webapps/traffic_router_${app}.war to /opt/traffic_router/webapps/${app}/ROOT.war"
		/bin/ln -s /opt/traffic_router/webapps/traffic_router_${app}.war /opt/traffic_router/webapps/${app}/ROOT.war
	fi
done

if [ -f /opt/traffic_router/conf/*.crt ]; then
	cd /opt/traffic_router/conf
	for file in *.crt; do
		alias=$(echo $file |sed -e 's/.crt//g' |tr [:upper:] [:lower:])
		cacerts=$(/bin/find $(dirname $(readlink -f $(which java)))/.. -name cacerts)
		keytool=$(dirname $(readlink -f $(which java)))/keytool
		$keytool -list -alias $alias -keystore $cacerts -storepass changeit -noprompt > /dev/null

		if [ $? -ne 0 ]; then
			echo "Installing certificate ${file}.."
			$keytool -import -trustcacerts -file $file -alias $alias -keystore $cacerts -storepass changeit -noprompt
		fi
	done
fi
