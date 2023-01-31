package version

import (
	vv "github.com/hashicorp/go-version"
)

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

// ShouldUpdate compares 2 versions (semver)
// return true if there is a major or 2 minor differences
func ShouldUpdate(current string, new string) bool {
	// don't force developers to update
	if current == "v0.0.0-dev" {
		return false
	}

	currentVer, err := vv.NewSemver(current)
	if err != nil {
		return false
	}

	newVer, err := vv.NewVersion(new)
	if err != nil {
		return false
	}

	// validate
	if len(newVer.Segments()) != 3 || len(currentVer.Segments()) != 3 {
		return false
	}

	// don't force updates for pre-releases
	if len(newVer.Prerelease()) > 0 {
		return false
	}

	// calculate diff in versions
	maj := newVer.Segments()[0] - currentVer.Segments()[0]
	min := newVer.Segments()[1] - currentVer.Segments()[1]

	// error case when current major version is newer than the one available for download (unlikely)
	if maj < 0 {
		return false
	}

	// if there is a major version difference, we force the update
	if maj >= 1 {
		return true
	}

	// if there are 2 minor versions difference, we force the update
	if min >= 2 {
		return true
	}

	return false
}
