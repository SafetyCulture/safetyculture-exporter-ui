#!/bin/bash

export PATH=${PATH}:`go env GOPATH`/bin
echo "building on AMD64"
wails build -platform darwin/amd64 -clean

cd build/bin
mkdir macOS

echo "Zipping Package"
ditto -c -k --keepParent ./build/bin/safetyculture-exporter.app ./build/bin/macOS/exporter-darwin-amd64.zip
echo "Cleaning up"
rm -rf ./build/bin/safetyculture-exporter.app

echo "building on ARM64"
wails build -platform darwin/arm64 -clean
echo "Zipping Package"
ditto -c -k --keepParent ./build/bin/safetyculture-exporter.app ./build/bin/macOS/exporter-darwin-arm64.zip
echo "Cleaning up"
rm -rf ./build/bin/safetyculture-exporter.app
