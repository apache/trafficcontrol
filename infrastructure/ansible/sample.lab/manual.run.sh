#!/usr/bin/env bash

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

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ ! -f ${DIR}/.vault_pass.txt ]]; then
  if [[ -z ${ANSIBLEPASS+x} ]]; then
    echo "Please enter your Ansible Vault Password: "
    read -sr ANSIBLEPASS_INPUT
    ANSIBLEPASS=$ANSIBLEPASS_INPUT
  fi
  echo $ANSIBLEPASS > "${DIR}/.vault_pass.txt"
fi

export ANSIBLE_ROLES_PATH="${DIR}/../roles"
# install any needed ansible-galaxy roles here
gilt overlay

pushd "${DIR}/out"
  find . -not -name '.dockerignore' -not -name '.gitignore' -not -name 'ssl' -not -path './ssl/lab.rootca.*' -not -path './ssl/lab.intermediateca.*' -delete
popd
pushd "${DIR}/inventory"
  find . -not -name '.dockerignore' -not -name '.gitignore' -delete
popd

# do whatever is necessary to kick off your provisioning and steady-state layers here

ansible-playbook -i "${DIR}/inventory" --vault-id "${DIR}/.vault_pass.txt" "${DIR}/../steady-state.yml"
ansible-playbook -i "${DIR}/inventory" --vault-id "${DIR}/.vault_pass.txt" "${DIR}/../sample.driver.playbook.yml"
ansible-playbook -i "${DIR}/inventory" --vault-id "${DIR}/.vault_pass.txt" "${DIR}/../influxdb_relay.yml"

rm -f "${DIR}/.vault_pass.txt"
