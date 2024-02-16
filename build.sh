#!/bin/bash
set -euxo pipefail


# chart version
number=$(cat version.txt)
number=$((number+1))
echo $number > version.txt
version="0.0.$number"
# use your own repo  
repository=michaelcourcy

# enter cmd directory to build images 
GOOS=linux GOARCH=amd64 go build
docker build --platform=linux/amd64 -t $repository/ceph-callback:$version-amd64 .
docker push $repository/ceph-callback:$version-amd64
rm ceph-callback
GOOS=linux GOARCH=arm64 go build
docker build --platform=linux/arm64 -t $repository/ceph-callback:$version-arm64 .
docker push $repository/ceph-callback:$version-arm64
rm ceph-callback
docker manifest create \
    $repository/ceph-callback:$version \
    $repository/ceph-callback:$version-amd64 \
    $repository/ceph-callback:$version-arm64
docker manifest push $repository/ceph-callback:$version
docker pull $repository/ceph-callback:$version
docker tag $repository/ceph-callback:$version $repository/ceph-callback:latest
docker push $repository/ceph-callback:latest