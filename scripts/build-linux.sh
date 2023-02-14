#! /bin/bash

echo "BUILDING WITH VERSION: ${VERSION}"

wails build -platform linux/amd64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version=${VERSION}"

mv ./build/bin/safetyculture-exporter .

mv safetyculture-exporter exporter-linux-amd64
