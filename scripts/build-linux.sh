#! /bin/bash

echo "DEBUGGING VARIABLES"
echo "TAG> ${ github.ref_name }"

wails build -platform linux/amd64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version=${github.ref_name}"

tar -czf exporter-linux-amd64.tar.gz ./build/bin/safetyculture-exporter
