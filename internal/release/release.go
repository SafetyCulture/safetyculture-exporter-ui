// Package release checks the SafetyCulture/mitti-exporter GitHub repository for
// the latest desktop app release. The final safetyculture-exporter-ui release
// uses it to migrate existing users to the new project home.
package release

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	osRuntime "runtime"
	"time"

	"github.com/hashicorp/go-version"
)

const (
	apiURL = "https://api.github.com/repos/SafetyCulture/mitti-exporter/releases/latest"

	// Desktop app assets published by the mitti-exporter release pipeline.
	// These are deliberately exact names: the same release also ships CLI
	// assets (e.g. mitti-exporter-darwin-arm64.zip, mitti-exporter-windows-x86_64.zip)
	// which must never be picked for a desktop app update.
	assetDarwinUniversal = "mitti-exporter-ui-darwin-universal.zip"
	assetWindows         = "mitti-exporter-ui-windows-x86_64.exe"

	requestTimeout = 10 * time.Second
)

// Info describes the latest available release of the mitti-exporter desktop app.
type Info struct {
	Version     string
	DownloadURL string
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

// Check returns the latest mitti-exporter desktop app release info for the
// current platform, or nil if it cannot be determined or the running version
// is not older than that release. It never returns an error: any failure
// (network, API, malformed response) results in nil so callers fail open and
// the app remains usable.
func Check(ctx context.Context, currentVersion string) *Info {
	return checkAt(ctx, apiURL, currentVersion)
}

func checkAt(ctx context.Context, url, currentVersion string) *Info {
	// A nil context must never disable the release check: the request falls
	// back to a background context so the migration still runs.
	if ctx == nil {
		ctx = context.Background()
	}

	client := &http.Client{Timeout: requestTimeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "safetyculture-exporter-ui")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil
	}

	var payload ghRelease
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil
	}

	latest := payload.TagName
	if latest == "" {
		return nil
	}

	if !isNewer(currentVersion, latest) {
		return nil
	}

	return &Info{
		Version:     latest,
		DownloadURL: assetURL(osRuntime.GOOS, payload.Assets),
	}
}

// isNewer reports whether latest is newer than currentVersion. When the
// running version cannot be parsed it reports true and lets the caller decide
// via version.ShouldUpdate.
func isNewer(currentVersion, latest string) bool {
	current, err := version.NewVersion(currentVersion)
	if err != nil {
		return true
	}
	next, err := version.NewVersion(latest)
	if err != nil {
		return false
	}
	return next.GreaterThan(current)
}

// assetURL returns the download URL of the desktop app asset for the given
// OS. It returns an empty string when the asset is missing, which callers
// treat as "release not downloadable" (no forced update).
func assetURL(goos string, assets []ghAsset) string {
	var want string
	switch goos {
	case "darwin":
		want = assetDarwinUniversal
	case "windows":
		want = assetWindows
	default:
		return ""
	}

	for _, asset := range assets {
		if asset.Name == want {
			return asset.BrowserDownloadURL
		}
	}
	return ""
}
