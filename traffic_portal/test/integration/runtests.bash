#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#
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