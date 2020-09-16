#!/usr/bin/env sh
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# shellcheck shell=ash
trap 'exit_code=$?; [ $exit_code -ne 0 ] && echo "Error on line ${LINENO} of ${0}"; cleanup; exit $exit_code' EXIT;
set -o errexit -o nounset -o pipefail;


# Fix ownership of output files
#  $1 is file or dir with correct ownership
#  remaining args are files/dirs to be fixed, recursively
setowner() {
	own=$(stat -c%u:%g "$1" 2>/dev/null || stat -f%u:%g "$1")
	shift
	[ -n "$*" ] && chown -R "${own}" "$@"
}

cleanup() {
	setowner "$tc_volume" "${tc_volume}/dist"
}

set -o xtrace;

if ! script_path="$(readlink "$0")"; then
	script_path="$0";
fi;
cd "$(dirname "$script_path")";
if ! tc_volume="$(git rev-parse --show-toplevel 2>/dev/null)"; then
	tc_volume='/trafficcontrol'; # Default repo location for ATC builder Docker images
fi;

# set owner of dist dir -- cleans up existing dist permissions...
export GOPATH=/tmp/go GOOS="${GOOS:-linux}";
tc_dir=${GOPATH}/src/github.com/apache/trafficcontrol;
if which cygpath 2>/dev/null; then
	GOPATH="$(cygpath -w "$GOPATH")" # cygwin compatibility
fi
(mkdir -p "$GOPATH"
 cd "$GOPATH"
 mkdir -p src pkg bin "$(dirname "$tc_dir")"
)
rsync -a --exclude=/dist --exclude=/.m2 "${tc_volume}/" "$tc_dir";
if ! [ -d ${tc_dir}/.git ]; then
	rsync -a "${tc_volume}/.git" $tc_dir; # Docker for Windows compatibility
fi

cd "$tc_dir"
# In case the mirrored repo already exists, remove gitignored files
git clean -fdX

rm -rf "dist"
mkdir -p "${tc_volume}/dist"
ln -sf "${tc_volume}/dist" "dist"

if [ $# -eq 0 ]; then
	set -- tarball traffic_monitor traffic_ops traffic_ops_ort traffic_portal traffic_router traffic_stats grove grove/grovetccfg docs
fi

for project in "$@"; do
	./build/build.sh "${project}" 2>&1 | tee "dist/build-${project//\//-}.log"
done
