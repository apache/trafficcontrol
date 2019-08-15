#!/usr/bin/env bash

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

#----------------------------------------
set -e
go get "github.com/wadey/gocovmerge"
touch result.txt
packages=( "$@" )
coverage_out=()
i=1
for pkg in ${packages[@]} ; do
    for d in $(go list $pkg | grep -v vendor); do
        file="$i.out"
        go test -v -coverprofile=$file $d | tee -a result.txt
        cat result.txt
        if [ -f $file ]; then
            coverage_out+=( $file )
        fi
        ((i++))
    done
done 
gocovmerge ${coverage_out[*]} > coverage.out
go tool cover -func=coverage.out
cat result.txt | go-junit-report --package-name=golang.test.tm --set-exit-code > /junit/golang.test.tm.xml 
cat /junit/golang.test.tm.xml
