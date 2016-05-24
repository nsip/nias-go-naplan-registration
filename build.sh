#!/bin/bash

mkdir -p build/Mac
mkdir -p build/Win32
mkdir -p build/Win64

echo "Building Mac binaries..."
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
cd ../csvxmlconverter
go get
go build
cd ../webui
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
cd ../csvxmlconverter
GOOS=windows GOARCH=amd64 go build
cd ../webui
GOOS=windows GOARCH=amd64 go build
echo "...all Windows64 binaries built..."
cd ..
echo "Building Windows32 binaries..."
cd ./aggregator
GOOS=windows GOARCH=386 go build -o aggregator32.exe
cd ../aslvalidator
GOOS=windows GOARCH=386 go build -o aslvalidator32.exe 
cd ../dobvalidator
GOOS=windows GOARCH=386 go build -o dobvalidator32.exe
cd ../idvalidator
GOOS=windows GOARCH=386 go build -o idvalidator32.exe
cd ../schemavalidator
GOOS=windows GOARCH=386 go build -o schemavalidator32.exe
cd ../csvxmlconverter
GOOS=windows GOARCH=386 go build -o csvxmlconverter32.exe
cd ../webui
GOOS=windows GOARCH=386 go build -o webui32.exe
echo "...all Windows32 binaries built..."
echo "go-nias Build Complete."





