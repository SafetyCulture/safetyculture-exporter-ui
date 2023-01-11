#! /bin/bash

echo "DEBUGGING VARIABLES"
echo "MICK DEBUG 1 ${{ github.ref_name }}"
echo "MICK DEBUG 2 ${{ github.ref }}"
echo "MICK DEBUG 3 ${{ env.Version }}"

wails build -platform linux/amd64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version={{.Version}}"

tar -czf exporter-linux-amd64.tar.gz ./build/bin/safetyculture-exporter
