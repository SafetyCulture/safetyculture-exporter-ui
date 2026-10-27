# SafetyCulture Exporter UI

> **Important:** This repository is no longer maintained and will be archived. The SafetyCulture Exporter desktop app is now published by the [mitti-exporter](https://github.com/SafetyCulture/mitti-exporter) project — download the latest release from <https://github.com/SafetyCulture/mitti-exporter/releases/latest>.
>
> The final release from this repository (v2.0.0) is a migration release: existing installs detect it through the built-in updater and are moved onto the mitti-exporter desktop app.

SafetyCulture Exporter UI is a desktop application that’s available to all our Premium and Enterprise customers. You can use the SafetyCulture Exporter UI to export your data, such as inspections, templates, schedules, and actions, to multiple formats that can be used for business intelligence tools or record keeping.

https://help.safetyculture.com/en_us/1064272525-rJ3PYPFEP

### How to run locally

Make sure you have installed wails cli: https://wails.io/docs/gettingstarted/installation

INSTALL WAILS: `go install github.com/wailsapp/wails/v2/cmd/wails@v2.9.1`

BUILD: `wails build -clean -debug`

RUN WEB: `wails dev`
