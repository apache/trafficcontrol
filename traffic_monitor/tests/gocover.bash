#!/usr/bin/env bash
set -e
go get "github.com/wadey/gocovmerge"
echo "" > result.txt
packages=( "$@" )
coverage_out=()
i=1
for pkg in ${packages[@]} ; do
    for d in $(go list $pkg | grep -v vendor); do
        file="$i.out"
        go test -v -coverprofile=$file $d >> result.txt
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
chmod 777 -R /junit && cat /junit/golang.test.tm.xml
