#!/bin/bash
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
keytool=$(dirname $(realpath "$(which java)"))/keytool
cd /opt/traffic_router/conf
if [ -f /opt/traffic_router/conf/*.crt ]; then
	for file in *.crt; do
		alias=$(echo $file |sed -e 's/.crt//g' |tr [:upper:] [:lower:])
		cacerts=$(/bin/find $(dirname $(realpath "$(which java)"))/.. -name cacerts)
		$keytool -list -alias $alias -keystore $cacerts -storepass changeit -noprompt > /dev/null

		if [ $? -ne 0 ]; then
			echo "Installing certificate ${file}.."
			$keytool -import -trustcacerts -file $file -alias $alias -keystore $cacerts -storepass changeit -noprompt
		fi
	done
fi


echo -e "
cd /opt/traffic_router/conf

keytool=\$(dirname \$(realpath \$(which java)))/keytool

if [ ! -f /opt/traffic_router/conf/keyStore.jks ]; then \n
    \$keytool -genkeypair -v -alias \$(hostname -f) -dname \"CN=\$(hostname -f), OU=APIDefault, O=Apache Traffic Control, L=Denver, ST=Colorado, C=US\" -keystore \$(pwd)/keyStore.jks -storepass changeit -keyalg RSA -ext KeyUsage=\"digitalSignature,keyEncipherment,keyCertSign\" -ext BasicConstraints:\"critical=ca:true\" -storetype JKS -validity 3650
    \$keytool -exportcert -v -alias \$(hostname -f) -file \$(hostname -f).crt -keypass changeit -storepass changeit -keystore \$(pwd)/keyStore.jks -rfc
fi" >> generatingCerts.sh
chmod 755 generatingCerts.sh

echo "Traffic Router installed successfully."

systemctl daemon-reload
echo "Start with 'systemctl start traffic_router'"

