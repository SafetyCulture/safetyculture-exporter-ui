package version

// This variable should be overridden at build time using ldflags.
var version string = "v0.0.0-dev"

const integrationID string = "safetyculture-exporter-ui"

// GetVersion returns current version of the app
func GetVersion() string {
	return version
}

// GetIntegrationID returns the integration id
func GetIntegrationID() string {
	return integrationID
}
