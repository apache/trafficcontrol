source /to-access.sh

curl -kvs -XPUT -H 'Content-Type:application/xml' "https://$TV_ADMIN_USER:$TV_ADMIN_PASSWORD@$TV_FQDN:$TV_HTTPS_PORT/search/schema/sslkeys" -d @/sslkeys.xml 

curl -kvs -XPUT -H 'Content-Type:application/json' "https://$TV_ADMIN_USER:$TV_ADMIN_PASSWORD@$TV_FQDN:$TV_HTTPS_PORT/search/index/sslkeys" -d '{"schema":"sslkeys"}'

curl -kvs -XPUT -H 'Content-Type:application/json' "https://$TV_ADMIN_USER:$TV_ADMIN_PASSWORD@$TV_FQDN:$TV_HTTPS_PORT/buckets/ssl/props" -d'{"props":{"search_index":"sslkeys"}}'
