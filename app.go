package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	osRuntime "runtime"
	"time"

	"github.com/SafetyCulture/safetyculture-exporter-ui/internal/version"
	exporterAPI "github.com/SafetyCulture/safetyculture-exporter/pkg/api"
	"github.com/SafetyCulture/safetyculture-exporter/pkg/httpapi"
	"github.com/SafetyCulture/safetyculture-exporter/pkg/update"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const gitRepoExporterUI string = "safetyculture-exporter-ui"

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

	ver := exporterAPI.AppVersion{
		IntegrationID:      version.GetIntegrationID(),
		IntegrationVersion: version.GetVersion(),
	}

	a.exporter, err = exporterAPI.NewSafetyCultureExporter(a.cm.Configuration, &ver)
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

// SelectDirectory opens a directory dialog and returns the path of the selected directory
func (a *App) SelectDirectory(currentDir string) string {
	var defaultDir string
	homeDir, err := os.UserHomeDir()
	if len(currentDir) == 0 {
		if err != nil {
			runtime.LogErrorf(a.ctx, "failed to get the working directory, %v", err)
			panic("failed to get working directory")
		}
		defaultDir = homeDir
	}

	defaultDir = currentDir

	directoryDialog, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:     defaultDir,
		CanCreateDirectories: true,
	})

	if err != nil {
		runtime.LogErrorf(a.ctx, "can't open directory dialog, %v", err)
		return homeDir
	}

	return directoryDialog
}

func (a *App) SelectSettingsDirectory() {
	settingsDir, err := GetSettingDirectoryPath()
	if err != nil {
		runtime.LogErrorf(a.ctx, "can't open directory dialog, %v", err)
		return
	}

	_, err = runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: settingsDir,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "*.log",
				Pattern:     "*.log",
			},
		},
	})

	if err != nil {
		runtime.LogErrorf(a.ctx, "can't open directory dialog, %v", err)
		return
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) ExportCSV() {
	a.exporter.RunCSV()
}

func (a *App) ExportSQL() {
	a.exporter.RunSQL()
}

func (a *App) GetTemplates() []exporterAPI.TemplateResponseItem {
	return a.exporter.GetTemplateList()
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

	cfg := httpapi.ClientCfg{
		Addr:                a.cm.Configuration.API.URL,
		AuthorizationHeader: fmt.Sprintf("Bearer %s", apiKey),
		IntegrationID:       version.GetIntegrationID(),
		IntegrationVersion:  version.GetVersion(),
	}
	c := httpapi.NewClient(&cfg, apiOpts...)
	res, err := httpapi.WhoAmI(a.ctx, c)

	if err != nil {
		runtime.LogErrorf(a.ctx, "cannot check WhoAmI: %s", err.Error())
		return false
	}

	if res != nil && (res.UserID == "" || res.OrganisationID == "") {
		runtime.LogErrorf(a.ctx, "cannot validate the credentials for the given ApiKey: %s", err.Error())
		return false
	}

	runtime.LogInfo(a.ctx, "saving the key")

	if apiKey != a.cm.Configuration.AccessToken {
		// save configuration
		a.cm.Configuration.AccessToken = apiKey
		if err := a.cm.SaveConfiguration(); err != nil {
			runtime.LogErrorf(a.ctx, "cannot save configuration: %s", err.Error())
		}

		ver := exporterAPI.AppVersion{
			IntegrationID:      version.GetIntegrationID(),
			IntegrationVersion: version.GetVersion(),
		}

		a.exporter, err = exporterAPI.NewSafetyCultureExporter(a.cm.Configuration, &ver)
		if err != nil {
			runtime.LogError(a.ctx, err.Error())
			panic("failed to re-load configuration")
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
	a.exporter.SetConfiguration(cfg)
	a.exporter.CleanExportStatus()
}

func (a *App) GetUserHomeDirectory() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		runtime.LogErrorf(a.ctx, "failed to find user's home directory, %v", err)
	}
	return dir
}

func (a *App) ReadExportStatus() {

	for {
		exportStatus := a.exporter.GetExportStatus()

		for _, item := range exportStatus.Feeds {
			runtime.EventsEmit(a.ctx, "update-"+item.FeedName, item)
		}

		if exportStatus.ExportStarted && exportStatus.ExportCompleted {
			runtime.EventsEmit(a.ctx, "finished-export", true)
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
}

type VersionResponse struct {
	Current      string `json:"current"`
	Latest       string `json:"latest"`
	DownloadURL  string `json:"download_url"`
	ShouldUpdate bool   `json:"should_update"`
}

func (a *App) ReadVersion() *VersionResponse {
	var current = version.GetVersion()
	var latest string
	var downloadURL string
	var shouldUpdate bool

	releaseInfo := update.Check(current, gitRepoExporterUI)
	if releaseInfo != nil {
		latest = releaseInfo.Version
		downloadURL = releaseInfo.DownloadURL
		shouldUpdate = version.ShouldUpdate(current, latest)
	}

	return &VersionResponse{
		Current:      current,
		Latest:       latest,
		DownloadURL:  downloadURL,
		ShouldUpdate: shouldUpdate,
	}
}

func (a *App) ReadBuild() string {
	return osRuntime.GOOS
}

func (a *App) CancelExport() {
	a.exporter.CancelExport()
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
