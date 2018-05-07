#!/bin/sh

if [[ ! -f ./records.config ]]; then
	echo "Error! Missing records.config file!" >&2
	exit 1
fi

if [[ ! -f ./remap.config ]]; then
	echo "Error! Missing remap.config file!" >&2
	exit 1
fi

mv ./records.config ./records.config.bak

NEWCONFIG=$(sed -re 's/CONFIG proxy\.config\.http\.cache\.http INT 0/CONFIG proxy\.config\.http\.cache\.http INT 1/' ./records.config.bak)
NEWCONFIG=$(echo $NEWCONFIG | sed -re 's/CONFIG proxy\.config\.reverse_proxy\.enabled INT 0/CONFIG proxy\.config\.reverse_proxy\.enabled INT 1/')
NEWCONFIG=$(echo $NEWCONFIG | sed -re 's/CONFIG proxy\.config\.url_remap\.remap_required INT 0/CONFIG proxy\.config\.url_remap\.remap_required INT 1/')
NEWCONFIG=$(echo $NEWCONFIG | sed -re 's/CONFIG proxy\.config\.url_remap\.pristine_host_hdr INT 0/CONFIG proxy\.config\.url_remap\.pristine_host_hdr INT 1/')
NEWCONFIG=$(echo $NEWCONFIG | sed -re 's/CONFIG proxy\.config\.http\.server_ports STRING ..*/CONFIG proxy\.config\.http\.server_ports STRING 8080 8080:ipv6/')

echo $NEWCONFIG > records.config

cp ./remap.config ./remap.config.bak
echo "map          http://proxied-origin/      http://origin/" >> remap.config
echo "reverse_map  http://origin/  http://proxied-origin/" >> remap.config

# echo "regex_map: http://(.*)/ http://origin/" >> remap.config
