#!/bin/bash
set -euxo pipefail

# use your own repo  
repository=michaelcourcy

# chart version
number=$(cat version.txt)
number=$((number+1))
echo $number > version.txt
version="0.0.$number"

# enter cmd directory to build images 
GOOS=linux GOARCH=amd64 go build
docker build --platform=linux/amd64 -t $repository/ceph-callback:$version-amd64 .
docker push $repository/ceph-callback:$version-amd64
rm ceph-callback