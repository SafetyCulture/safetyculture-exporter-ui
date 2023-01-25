package version

import "github.com/SafetyCulture/safetyculture-exporter/pkg/update"

// This variable should be overridden at build time using ldflags.
var version string = "v0.0.0-dev"

const integrationID string = "safetyculture-exporter-ui"
const gitRepoExporterUI string = "safetyculture-exporter-ui"

// GetVersion returns current version of the app
func GetVersion() string {
	return version
}

// GetIntegrationID returns the integration id
func GetIntegrationID() string {
	return integrationID
}

// GetLatestVersion will connect to GitHub and get the latest tag if newer
func GetLatestVersion() string {
	releaseInfo := update.Check(version, gitRepoExporterUI)
	if releaseInfo == nil {
		return "error"
	}
	return releaseInfo.Version
}
