#!/bin/bash

wails build -platform windows/amd64 -clean

cd build/bin
tar -czf exporter-windows-amd64.tar.gz safetyculture-exporter.exe
