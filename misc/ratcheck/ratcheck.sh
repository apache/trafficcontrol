# ----------------------------------------------------------------------
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
# ----------------------------------------------------------------------
#                      incubator-trafficcontrol-rat
#
#          This file captures the Apache Jenkins build script
#                    Copied from HAWQ-rat build job
# ----------------------------------------------------------------------
OPTIND=1

#ABSPATH=$(pushd "$(dirname "$0")" > /dev/null; pwd; popd > /dev/null)
#cd ${ABSPATH}/../..

# Use apache-rat-0.13 which has support for .go files
ratver=apache-rat-0.13-20171021.191031-67.jar
ratdestname=apache-rat-0.13.SNAPSHOT.jar
ratdir=$(mktemp -ud)
targetdir=$(pwd)
outreport=$targetdir

function get_abs_path(){
  mkdir -p ${1}
  ABSPATH=$(pushd "${1}" > /dev/null; pwd; popd > /dev/null)
  echo $ABSPATH
}

function show_help {
  echo "A script to download and run the RAT license checker"
  echo "Usage: [-h] [-p x] [-s x] [-d x] [-t x]"
  echo "  -h  show this help text"
  echo "  -p  location used to download and execute RAT from (default: ${ratdir})"
  echo "  -s  source RAT jarfile name (default: ${ratver})"
  echo "  -d  destination RAT jarfile name (default: ${ratdestname})"
  echo "  -t  target directory to run RAT on (default: ${targetdir})"
  echo "  -o  output directory for ratreport.txt (default: ${outreport})"
}

while getopts "h?p:s:d:t:o:" opt; do
  case "$opt" in
  h|\?)
    show_help
    exit 0
    ;;
  p)
    ratdir=$(get_abs_path $OPTARG)
    ;;
  s)
    ratver=$OPTARG
    ;;
  d)
    ratdestname=$OPTARG
    ;;
  t)
    targetdir=$(get_abs_path $OPTARG)
    ;;
  t)
    outreport=$(get_abs_path $OPTARG)
    ;;
  esac
done
shift $((OPTIND-1))
ratjar="$ratdir/$ratdestname"

set -exu
# Check if NOTICE file year is current
grep "Copyright $(date +"%Y") The Apache Software Foundation" "${targetdir}/NOTICE"
set +x

badfile_extentions="class jar tar tgz zip"
badfiles_found=false

for extension in ${badfile_extentions}; do
    echo "Searching for ${extension} files:"
    badfile_count=$(find ${targetdir} -name "*.${extension}" | wc -l)
    if [ ${badfile_count} != 0 ]; then
        echo "----------------------------------------------------------------------"
        echo "FATAL: ${extension} files should not exist"
        echo "For ASF compatibility: the source tree should not contain"
        echo "binary (jar) files as users have a hard time verifying their"
        echo "contents."

        find ${targetdir} -name "*.${extension}"
        echo "----------------------------------------------------------------------"
        badfiles_found=true
    else
        echo "PASSED: No ${extension} files found."
    fi
done

if [ ${badfiles_found} = "true" ]; then
    exit 1
fi

set -x

curl -L -o $ratjar \
  https://repository.apache.org/content/repositories/snapshots/org/apache/rat/apache-rat/0.13-SNAPSHOT/${ratver}
curl -L -o $ratjar.sha1 \
  https://repository.apache.org/content/repositories/snapshots/org/apache/rat/apache-rat/0.13-SNAPSHOT/${ratver}.sha1

# Check sha1 on downloaded .jar
[[ $(sha1sum "${ratjar}" | awk '{print $1}') == $(cat "${ratjar}.sha1") ]] || \
   (echo "SHA1 check failed -- aborting!"; exit 1 )


# Run rat and generate report
java -jar "${ratjar}" -E "${targetdir}/.rat-excludes" -d "$targetdir" > "$outreport/ratreport.txt"


unknown=$(perl -lne 'print $1 if /(\d+) Unknown Licenses/' $outreport/ratreport.txt)


if [[ $unknown != 0 ]]; then
    echo "$unknown Unknown Licenses"
    perl -lne 'print if /Files with unapproved licenses:/ .. /^\*\*\*/' "$outreport/ratreport.txt" | \
        sed "s:${targetdir}/::"
    exit 1
fi
