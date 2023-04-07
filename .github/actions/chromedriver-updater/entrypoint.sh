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

trap 'echo "Error on line ${LINENO} of ${0}"; exit 1' ERR;
set -o errexit -o nounset

if [[ -z ${PROJECTS} ]]; then
  echo "PROJECTS environment variable not set!"
  exit 1
fi

python3 -m chromedriver_updater

projects=($(echo $PROJECTS | tr ',' ' '))
touch updates.txt

for proj in "${projects[@]}"
do
  package="./$proj"package.json
  if [[ ! -f "$package" ]]; then
    echo "Unable to find package.json in project directory $proj"
    continue
  fi

  pushd "./$proj" > /dev/null

  npm ci

  outdated=$(npm outdated | grep "chromedriver" || echo "" )

  if [[ -z $outdated ]]; then
    echo "$proj is up to date"
    popd > /dev/null
    continue
  fi

  latest=$(echo $outdated | awk '{print $4}' )
  wanted=$(echo $outdated | awk '{print $3}' )

  npm i --save-dev "chromedriver@$latest" --ignore-scripts > /dev/null

  popd > /dev/null
  echo -e "$proj:$wanted,$latest\n" >> updates.txt
done

python3 -m chromedriver_updater updates.txt
