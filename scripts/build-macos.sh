#!/bin/bash

echo "BUILDING WITH VERSION: ${VERSION}"

export PATH=${PATH}:`go env GOPATH`/bin
echo "building on AMD64"
wails build -platform darwin/amd64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version=${VERSION}"

echo "Zipping Package"
ditto -c -k --rsrc --keepParent ./build/bin/safetyculture-exporter.app ./exporter-darwin-amd64.zip
echo "Cleaning up"
rm -rf build/bin/SafetyCulture Exporter.app

echo "building on ARM64"
wails build -platform darwin/arm64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version=${VERSION}"

echo "Zipping Package"
ditto -c -k --rsrc --keepParent ./build/bin/safetyculture-exporter.app ./exporter-darwin-arm64.zip
echo "Cleaning up"
rm -rf ./build/bin/SafetyCulture Exporter.app
