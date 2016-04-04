#!/bin/bash

echo "Building Mac binaries..."
go get github.com/nats-io/gnatsd
# TODO - this will build in current directory change...
cd ../../nats-io/gnatsd
go build
cd -
cd ./aggregator
go get
go build 
cd ../aslvalidator
go get
go build
cd ../dobvalidator
go get
go build
cd ../idvalidator
go get
go build
cd ../schemavalidator
go get
go build
echo "...all Mac binaries built..."
cd ..
echo "Building Windows64 binaries..."
cd ./aggregator
GOOS=windows GOARCH=amd64 go build 
cd ../aslvalidator
GOOS=windows GOARCH=amd64 go build
cd ../dobvalidator
GOOS=windows GOARCH=amd64 go build
cd ../idvalidator
GOOS=windows GOARCH=amd64 go build
cd ../schemavalidator
GOOS=windows GOARCH=amd64 go build
echo "...all Windows binaries built..."
echo "go-nias Build Complete."
