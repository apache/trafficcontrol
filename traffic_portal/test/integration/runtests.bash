#!/bin/bash

# isntall dependencies
npm install 
webdriver-manager update 
nohup webdriver-manager start &

while ! curl -Lvsk "http://localhost:4444" 2>/dev/null >/dev/null; do
   echo "waiting for selenium server to start on 'http://localhost:4444'"
   sleep 1
done

tsc
protractor ./GeneratedCode/config.js --params.baseUrl=$1
rc=$?

exit $rc