#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#set -e

mysqld --user=root --initialize-insecure
temp=$(cat /var/log/mysql/mysqld.log | grep temp)
mysqld --user=root &


mkdir -p /var/www
cd /var/www
wget https://builds.matomo.org/matomo.zip
unzip matomo.zip > /dev/null

set-dns.sh
insert-self-into-dns.sh

until [[ -f "$X509_CA_ENV_FILE" ]]
do
  echo "Waiting on Shared SSL certificate generation"
  sleep 3
done

# Source the CIAB-CA shared SSL environment
until [[ -n "$X509_GENERATION_COMPLETE" ]]
do
  echo "Waiting on X509 vars to be defined"
  sleep 1
  source "$X509_CA_ENV_FILE"
done

#set-dns.sh
#ls /shared/dns
#insert-self-into-dns.sh

# Trust the CIAB-CA at the System level
cp $X509_CA_CERT_FULL_CHAIN_FILE /etc/pki/ca-trust/source/anchors
update-ca-trust extract

# Configuration of Traffic Portal
key=$X509_INFRA_KEY_FILE
cert=$X509_INFRA_CERT_FILE

sed -i 's@\!key@'"$key"'@g' /matomo.conf
sed -i 's@\!cert@'"$cert"'@g' /matomo.conf
cat /matomo.conf | grep ssl

mkdir -p /etc/nginx/sites-available/ /run/php-fpm
mv /matomo.conf /etc/nginx/conf.d/
mv /config.ini.php /var/www/matomo/config

mysql < /mysql.sql
mysql --database=matomo --user=matomo --password=twelve < /matomo.sql
php-fpm
nginx

chown -R apache:apache /var/www/matomo

echo "done"
tail -f /dev/null
