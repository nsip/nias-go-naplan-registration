#!/bin/bash

echo "Building Mac binaries..."
cd ./aggregator
go build 
cd ../aslvalidator
go build
cd ../dobvalidator
go build
cd ../idvalidator
go build
cd ../schemavalidator
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