package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	osRuntime "runtime"
	"strings"

	exporterAPI "github.com/SafetyCulture/safetyculture-exporter/pkg/api"
	"github.com/SafetyCulture/safetyculture-exporter/pkg/httpapi"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx      context.Context
	cm       *exporterAPI.ConfigurationManager
	exporter *exporterAPI.SafetyCultureExporter
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so, we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	settingsDir, err := GetSettingDirectoryPath()
	if err != nil {
		runtime.LogError(ctx, "failed to get settings directory")
		panic("failed to get settings directory")
	}

	if !checkForConfigFile(settingsDir) {
		runtime.LogInfof(ctx, "creating configuration file: %s/safetyculture-exporter.yaml", settingsDir)
		cm := exporterAPI.NewConfigurationManager(settingsDir, "safetyculture-exporter.yaml")
		if err := cm.SaveConfiguration(); err != nil {
			runtime.LogError(ctx, err.Error())
			panic("failed to save configuration")
		}
		a.cm = cm
	} else {
		runtime.LogInfof(ctx, "loading configuration file: %s/safetyculture-exporter.yaml", settingsDir)
		cm, err := exporterAPI.NewConfigurationManagerFromFile(settingsDir, "safetyculture-exporter.yaml")
		if err != nil {
			runtime.LogError(ctx, err.Error())
			panic("failed to load configuration")
		}
		a.cm = cm
	}

	a.exporter, err = exporterAPI.NewSafetyCultureExporterInferredApiClient(a.cm.Configuration)
	if err != nil {
		runtime.LogError(ctx, err.Error())
		panic("failed to load configuration")
	}
	a.ctx = ctx
}

func checkForConfigFile(basePath string) bool {
	if _, err := os.Stat(path.Join(basePath, "safetyculture-exporter.yaml")); os.IsNotExist(err) {
		return false
	}
	return true
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) ExportCSV() {

}

func (a *App) GetTemplates() []exporterAPI.TemplateResponseItem {
	res, err := a.exporter.GetTemplateList()
	if err != nil {
		runtime.LogErrorf(a.ctx, "cannot obtain the template list: %s", err.Error())
		return []exporterAPI.TemplateResponseItem{}
	}
	return res
}

// CheckApiKey validates the api key from the config file if it exists
func (a *App) CheckApiKey() bool {
	token := a.cm.Configuration.AccessToken
	if len(token) == 0 {
		return false
	}

	return a.ValidateApiKey(token)
}

// ValidateApiKey validates the api, returns true if valid, false otherwise
func (a *App) ValidateApiKey(apiKey string) bool {
	var apiOpts []httpapi.Opt

	c := httpapi.NewClient(a.cm.Configuration.API.URL, fmt.Sprintf("Bearer %s", apiKey), apiOpts...)
	res, err := c.WhoAmI(a.ctx)

	if err != nil {
		runtime.LogErrorf(a.ctx, "cannot check WhoAmI: %s", err.Error())
		return false
	}

	if res != nil && (res.UserID == "" || res.OrganisationID == "") {
		runtime.LogErrorf(a.ctx, "cannot validate the credentials for the given ApiKey: %s", err.Error())
		return false
	}

	runtime.LogInfo(a.ctx, "saving the key")

	if strings.Compare(apiKey, a.cm.Configuration.AccessToken) != 0 {
		// save configuration
		a.cm.Configuration.AccessToken = apiKey
		if err := a.cm.SaveConfiguration(); err != nil {
			runtime.LogErrorf(a.ctx, "cannot save configuration: %s", err.Error())
		}

		// reload configuration
		a.exporter, err = exporterAPI.NewSafetyCultureExporterInferredApiClient(a.cm.Configuration)
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
			panic("failed to load configuration")
		}
	}
	return true
}

// GetSettings gets the configuration from the YAML file
func (a *App) GetSettings() *exporterAPI.ExporterConfiguration {
	return a.cm.Configuration
}

// SaveSettings saves the configuration to the YAML file
func (a *App) SaveSettings(cfg *exporterAPI.ExporterConfiguration) {
	a.cm.Configuration = cfg
	if err := a.cm.SaveConfiguration(); err != nil {
		runtime.LogErrorf(a.ctx, "cannot save configuration: %s", err.Error())
	}
}

func CreateSettingsDirectory() (string, error) {
	settingDir, err := GetSettingDirectoryPath()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(settingDir); os.IsNotExist(err) {
		err := os.MkdirAll(settingDir, 0700)
		if err != nil {
			return "", errors.New("can't create settings directory")
		}
	}

	return settingDir, nil
}

func GetSettingDirectoryPath() (string, error) {
	switch osRuntime.GOOS {
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("can't get user's home directory")
		}
		return filepath.Join(homeDir, "/Library/Application Support/safetyculture-exporter"), nil
	default:
		wd, err := os.Getwd()
		if err != nil {
			return "", errors.New("can't get user's home directory")
		}
		return wd, nil
	}
}