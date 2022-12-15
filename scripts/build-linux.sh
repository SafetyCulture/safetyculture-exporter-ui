#! /bin/bash

wails build -platform linux/amd64 -clean

cd build/bin
tar -czf exporter-linux-amd64.tar.gz safetyculture-exporter
