#! /bin/bash

wails build -platform linux/amd64 -clean

cd build/bin
tar -czf safetyculture-exporter.tar.gz safetyculture-exporter
mkdir linux
mv safetyculture-exporter.tar.gz linux/