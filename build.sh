#!/bin/bash

set -e

mkdir -p build/Mac/go-nias
mkdir -p build/Win32/go-nias
mkdir -p build/Win64/go-nias
mkdir -p build/Linux64/go-nias
mkdir -p build/Linux32/go-nias
CWD=`pwd`

# MAC OS X (64 only)
echo "Building Mac binaries..."
go get github.com/nats-io/gnatsd
cd ../../nats-io/gnatsd
go build -o $CWD/build/Mac/go-nias/gnatsd
cd $CWD
cd ./aggregator
go get
go build -o $CWD/build/Mac/go-nias/aggregator
cd ../aslvalidator
go get
go build -o $CWD/build/Mac/go-nias/aslvalidator
cd ../dobvalidator
go get
go build -o $CWD/build/Mac/go-nias/dobvalidator
cd ../idvalidator
go get
go build -o $CWD/build/Mac/go-nias/idvalidator
cd ../schemavalidator
go get
go build -o $CWD/build/Mac/go-nias/schemavalidator
cd ../csvxmlconverter
go get
go build -o $CWD/build/Mac/go-nias/csvxmlconverter
cd ../webui
go get
go build -o $CWD/build/Mac/go-nias/webui
cd ..
rsync -a csvxmlconverter/templates webui/public schemavalidator/schemas aslvalidator/schoolslist test_data $CWD/build/Mac/go-nias/
cp gonias.sh $CWD/build/Mac/go-nias/
cp shutdown.sh $CWD/build/Mac/go-nias/
cd build/Mac
zip -r ../go-nias-Mac.zip go-nias/
cd $CWD
echo "...all Mac binaries built..."

# WINDOWS 64
echo "Building Windows64 binaries..."
cd ../../nats-io/gnatsd
go get
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/go-nias/gnatsd.exe
cd $CWD
cd ./aggregator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/go-nias/aggregator.exe
cd ../aslvalidator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/go-nias/aslvalidator.exe
cd ../dobvalidator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/go-nias/dobvalidator.exe
cd ../idvalidator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/go-nias/idvalidator.exe
cd ../schemavalidator
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/go-nias/schemavalidator.exe
cd ../csvxmlconverter
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/go-nias/csvxmlconverter.exe
cd ../webui
GOOS=windows GOARCH=amd64 go build -o $CWD/build/Win64/go-nias/webui.exe
cd ..
rsync -a csvxmlconverter/templates webui/public schemavalidator/schemas aslvalidator/schoolslist test_data $CWD/build/Win64/go-nias/
cp gonias.bat $CWD/build/Win64/go-nias/
cd build/Win64
zip -r ../go-nias-Win64.zip go-nias/
cd $CWD
echo "...all Windows64 binaries built..."

# WINDOWS 32
echo "Building Windows32 binaries..."
cd ../../nats-io/gnatsd
go get
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/go-nias/gnatsd.exe
cd $CWD
cd ./aggregator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/go-nias/aggregator.exe
cd ../aslvalidator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/go-nias/aslvalidator.exe 
cd ../dobvalidator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/go-nias/dobvalidator.exe
cd ../idvalidator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/go-nias/idvalidator.exe
cd ../schemavalidator
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/go-nias/schemavalidator.exe
cd ../csvxmlconverter
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/go-nias/csvxmlconverter.exe
cd ../webui
GOOS=windows GOARCH=386 go build -o $CWD/build/Win32/go-nias/webui.exe
cd ..
rsync -a csvxmlconverter/templates webui/public schemavalidator/schemas aslvalidator/schoolslist test_data $CWD/build/Win32/go-nias/
cp gonias.bat $CWD/build/Win32/go-nias/
cd build/Win32
zip -r ../go-nias-Win32.zip go-nias/
cd $CWD
echo "...all Windows32 binaries built..."
echo "go-nias Build Complete."


# LINUX (64)
echo "Building Linux64 binaries..."
cd ../../nats-io/gnatsd
go get
GOOS=linux GOARCH=amd64 go build -o $CWD/build/Linux64/go-nias/gnatsd
cd $CWD
cd ./aggregator
go get
GOOS=linux GOARCH=amd64 go build -o $CWD/build/Linux64/go-nias/aggregator
cd ../aslvalidator
go get
GOOS=linux GOARCH=amd64 go build -o $CWD/build/Linux64/go-nias/aslvalidator
cd ../dobvalidator
go get
GOOS=linux GOARCH=amd64 go build -o $CWD/build/Linux64/go-nias/dobvalidator
cd ../idvalidator
go get
GOOS=linux GOARCH=amd64 go build -o $CWD/build/Linux64/go-nias/idvalidator
cd ../schemavalidator
go get
GOOS=linux GOARCH=amd64 go build -o $CWD/build/Linux64/go-nias/schemavalidator
cd ../csvxmlconverter
go get
GOOS=linux GOARCH=amd64 go build -o $CWD/build/Linux64/go-nias/csvxmlconverter
cd ../webui
go get
GOOS=linux GOARCH=amd64 go build -o $CWD/build/Linux64/go-nias/webui
cd ..
rsync -a csvxmlconverter/templates webui/public schemavalidator/schemas aslvalidator/schoolslist test_data $CWD/build/Linux64/go-nias/
cp gonias.sh $CWD/build/Linux64/go-nias/
cp shutdown.sh $CWD/build/Linux64/go-nias/
cd build/Linux64
zip -r ../go-nias-Linux64.zip go-nias/
cd $CWD
echo "...all Linux64 binaries built..."

# LINUX (32)
echo "Building Linux32 binaries..."
cd ../../nats-io/gnatsd
go get
GOOS=linux GOARCH=386 go build -o $CWD/build/Linux32/go-nias/gnatsd
cd $CWD
cd ./aggregator
go get
GOOS=linux GOARCH=386 go build -o $CWD/build/Linux32/go-nias/aggregator
cd ../aslvalidator
go get
GOOS=linux GOARCH=386 go build -o $CWD/build/Linux32/go-nias/aslvalidator
cd ../dobvalidator
go get
GOOS=linux GOARCH=386 go build -o $CWD/build/Linux32/go-nias/dobvalidator
cd ../idvalidator
go get
GOOS=linux GOARCH=386 go build -o $CWD/build/Linux32/go-nias/idvalidator
cd ../schemavalidator
go get
GOOS=linux GOARCH=386 go build -o $CWD/build/Linux32/go-nias/schemavalidator
cd ../csvxmlconverter
go get
GOOS=linux GOARCH=386 go build -o $CWD/build/Linux32/go-nias/csvxmlconverter
cd ../webui
go get
GOOS=linux GOARCH=386 go build -o $CWD/build/Linux32/go-nias/webui
cd ..
rsync -a csvxmlconverter/templates webui/public schemavalidator/schemas aslvalidator/schoolslist test_data $CWD/build/Linux32/go-nias/
cp gonias.sh $CWD/build/Linux32/go-nias/
cp shutdown.sh $CWD/build/Linux32/go-nias/
cd build/Linux32
zip -r ../go-nias-Linux32.zip go-nias/
cd $CWD
echo "...all Linux32 binaries built..."

# LINUX (arm7)
echo "Building LinuxArm7 binaries..."
cd ../../nats-io/gnatsd
go get
GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/gnatsd
cd $CWD
cd ./aggregator
go get
GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/aggregator
cd ../aslvalidator
go get
GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/aslvalidator
cd ../dobvalidator
go get
GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/dobvalidator
cd ../idvalidator
go get
GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/idvalidator
cd ../schemavalidator
go get
GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/schemavalidator
cd ../csvxmlconverter
go get
GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/csvxmlconverter
cd ../webui
go get
GOOS=linux GOARCH=arm GOARM=7 go build -o $CWD/build/LinuxArm7/go-nias/webui
cd ..
rsync -a csvxmlconverter/templates webui/public schemavalidator/schemas aslvalidator/schoolslist test_data $CWD/build/LinuxArm7/go-nias/
cp gonias.sh $CWD/build/LinuxArm7/go-nias/
cp shutdown.sh $CWD/build/LinuxArm7/go-nias/
cd build/LinuxArm7
zip -r ../go-nias-LinuxArm7.zip go-nias/
cd $CWD
echo "...all LinuxArm7 binaries built..."

