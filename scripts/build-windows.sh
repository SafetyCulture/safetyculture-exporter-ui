#!/bin/bash

echo "BUILDING WITH VERSION: ${VERSION}"

wails build -platform windows/amd64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version=${VERSION}"

tar -czf exporter-windows-amd64.tar.gz ./build/bin/safetyculture-exporter.exe
