package main

import (
	"embed"
	"fmt"
	"github.com/SafetyCulture/safetyculture-exporter-ui/internal/version"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	extLogger "github.com/SafetyCulture/safetyculture-exporter/pkg/logger"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()

	settingsDir, err := CreateSettingsDirectory()
	if err != nil {
		panic(fmt.Sprintf("failed to get settings directory: %s", err.Error()))
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:         "SafetyCulture Exporter",
		Width:         1080,
		Height:        780,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
		LogLevel: logger.INFO,
		Logger:   extLogger.GetExporterLogger(settingsDir),
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Mac: &mac.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "SafetyCulture Exporter",
				Message: fmt.Sprintf("Version %v\n\nCopyright \u00a9 2023\nSafetyCulture Pty Ltd.", version.GetVersion()),
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
