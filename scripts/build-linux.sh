#! /bin/bash

echo "DEBUGGING VARIABLES"
echo ${{ github.github.actor }}
echo ${{ github.ref_name }}
echo ${{ github.ref }}
echo ${{ env.Version }}
echo ${{ env.TAG_NAME }}
echo ${{ github.repository }}

wails build -platform linux/amd64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version=${{ env.TAG_NAME }}"

tar -czf exporter-linux-amd64.tar.gz ./build/bin/safetyculture-exporter
