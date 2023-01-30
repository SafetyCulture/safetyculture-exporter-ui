#!/bin/bash

echo "BUILDING WITH VERSION: ${VERSION}"
echo "CGO ENABLED: ${CGO_ENABLED}"

wails build -platform windows/amd64 -clean -ldflags "-s -w -X github.com/SafetyCulture/safetyculture-exporter-ui/internal/version.version=${VERSION}" -skipbindings

zip -r exporter-windows-amd64.zip ./build/bin/safetyculture-exporter.exe
