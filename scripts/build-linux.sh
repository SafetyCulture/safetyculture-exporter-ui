#! /bin/bash

wails build -platform linux/amd64 -clean

cd build/bin
tar -czf exporter.tar.gz safetyculture-exporter
mkdir linux
mv exporter.tar.gz linux/
