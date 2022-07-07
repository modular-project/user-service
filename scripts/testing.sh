#!/bin/sh

for d in $(go list ./...); do
	echo "Testing the package $d"
	go test -v $d
done
