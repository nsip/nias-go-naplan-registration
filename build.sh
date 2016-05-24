#!/bin/bash

mkdir -p build/Mac
mkdir -p build/Win32
mkdir -p build/Win64
CWD=`pwd`

echo "Building Mac binaries..."
cd ../../nats-io/gnatsd
go get
go build -o $CWD/build/Mac/gnatsd
cd $CWD
cd ./aggregator
go get
go build -o $CWD/build/Mac/aggregator
cd ../aslvalidator
go get
go build -o $CWD/build/Mac/aslvalidator
cd ../dobvalidator
go get
go build -o $CWD/build/Mac/dobvalidator
cd ../idvalidator
go get
go build -o $CWD/build/Mac/idvalidator
cd ../schemavalidator
go get
go build -o $CWD/build/Mac/schemavalidator
cd ../csvxmlconverter
go get
go build -o $CWD/build/Mac/csvxmlconverter
cd ../webui
go get
go build -o $CWD/build/Mac/webui
rsync -a webui/public schemavalidator/schemas aslvalidator/schoolslist test_data $CWD/build/Mac/
cp gonias.sh $CWD/build/Mac/
echo "...all Mac binaries built..."
cd ..
echo "Building Windows64 binaries..."
cd ../../nats-io/gnatsd
go get
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/gnatsd.exe
cd $CWD
cd ./aggregator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/aggregator.exe
cd ../aslvalidator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/aslvalidator.exe
cd ../dobvalidator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/dobvalidator.exe
cd ../idvalidator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/idvalidator.exe
cd ../schemavalidator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/schemavalidator.exe
cd ../csvxmlconverter
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/csvxmlconverter.exe
cd ../webui
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/webui.exe
rsync -a webui/public schemavalidator/schemas aslvalidator/schoolslist test_data $CWD/build/Win64/
cp gonias.bat $CWD/build/Win64/
echo "...all Windows64 binaries built..."
cd ..
echo "Building Windows32 binaries..."
cd ../../nats-io/gnatsd
go get
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/gnatsd.exe
cd $CWD
cd ./aggregator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/aggregator.exe
cd ../aslvalidator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/aslvalidator.exe 
cd ../dobvalidator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/dobvalidator.exe
cd ../idvalidator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/idvalidator.exe
cd ../schemavalidator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/schemavalidator.exe
cd ../csvxmlconverter
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/csvxmlconverter.exe
cd ../webui
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/webui.exe
rsync -a webui/public schemavalidator/schemas aslvalidator/schoolslist test_data $CWD/build/Win32/
cp gonias.bat $CWD/build/Win32/
echo "...all Windows32 binaries built..."
echo "go-nias Build Complete."

