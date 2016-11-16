#!/bin/bash -x

output=$1
[[ -n $output ]] && output="-o $output"


cookiejar=/tmp/cookiejar
cred=/tmp/cred.json

cat >$cred <<-CREDS
	{ "u" : "$TO_USER", "p" : "$TO_PASSWORD" }
CREDS

curl -k -H "Accept: application/json" --cookie "$cookiejar" --cookie-jar "$cookiejar" -X POST --data @"$cred" "$TO_SERVER/api/1.2/user/login"
curl $output -k -s --cookie "$cookiejar" -X GET "$TO_SERVER/dbdump"
