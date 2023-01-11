#! /bin/bash

echo "DEBUGGING VARIABLES"
echo "MICK DEBUG 1 ${{ github.github.actor }}"
echo "MICK DEBUG 2 ${{ github.ref_name }}"
echo "MICK DEBUG 3 ${{ github.ref }}"
echo "MICK DEBUG 4 ${{ env.Version }}"
echo "MICK DEBUG 5 ${{ env.TAG_NAME }}"

wails build -platform linux/amd64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version=${{ env.TAG_NAME }}"

tar -czf exporter-linux-amd64.tar.gz ./build/bin/safetyculture-exporter
