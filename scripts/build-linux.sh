#! /bin/bash

echo "BUILDING WITH VERSION: ${VERSION}"

wails build -platform linux/amd64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version=${VERSION}"

tar -czf exporter-linux-amd64.tar.gz ./build/bin/safetyculture-exporter
