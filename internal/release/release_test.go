package release

import (
	"net/http"
	"net/http/httptest"
	osRuntime "runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func goosForTest() string { return osRuntime.GOOS }

// mittiV630Assets mirrors the assets of the real mitti-exporter v6.3.0
// release, in API order, including the CLI assets that must never be selected.
var mittiV630Assets = []ghAsset{
	{Name: "mitti-exporter-darwin-amd64.zip", BrowserDownloadURL: "https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-darwin-amd64.zip"},
	{Name: "mitti-exporter-darwin-arm64.zip", BrowserDownloadURL: "https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-darwin-arm64.zip"},
	{Name: "mitti-exporter-linux-amd64.zip", BrowserDownloadURL: "https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-linux-amd64.zip"},
	{Name: "mitti-exporter-ui-darwin-universal.zip", BrowserDownloadURL: "https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-ui-darwin-universal.zip"},
	{Name: "mitti-exporter-ui-windows-x86_64.exe", BrowserDownloadURL: "https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-ui-windows-x86_64.exe"},
	{Name: "mitti-exporter-windows-x86_64.zip", BrowserDownloadURL: "https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-windows-x86_64.zip"},
}

func newReleaseServer(t *testing.T, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
}

func TestCheckAt_MigrationReleaseSeesMittiAsNewer(t *testing.T) {
	srv := newReleaseServer(t, `{"tag_name":"v6.3.0","assets":[
		{"name":"mitti-exporter-darwin-arm64.zip","browser_download_url":"https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-darwin-arm64.zip"},
		{"name":"mitti-exporter-ui-darwin-universal.zip","browser_download_url":"https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-ui-darwin-universal.zip"},
		{"name":"mitti-exporter-ui-windows-x86_64.exe","browser_download_url":"https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-ui-windows-x86_64.exe"},
		{"name":"mitti-exporter-windows-x86_64.zip","browser_download_url":"https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-windows-x86_64.zip"}]}`)
	defer srv.Close()

	info := checkAt(srv.URL, "v2.0.0")
	require.NotNil(t, info)
	assert.EqualValues(t, "v6.3.0", info.Version)
	// the download URL for the host platform matches the expected desktop
	// asset (darwin/windows), or is empty on unsupported platforms
	assert.EqualValues(t, assetURL(goosForTest(), mittiV630Assets), info.DownloadURL)
}

func TestCheckAt_ReturnsNilWhenCurrentIsNewerOrEqual(t *testing.T) {
	srv := newReleaseServer(t, `{"tag_name":"v6.3.0","assets":[]}`)
	defer srv.Close()

	assert.Nil(t, checkAt(srv.URL, "v6.3.0"))
	assert.Nil(t, checkAt(srv.URL, "v6.4.0"))
}

func TestCheckAt_ReturnsNilOnHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	assert.Nil(t, checkAt(srv.URL, "v2.0.0"))
}

func TestCheckAt_ReturnsNilOnMalformedResponse(t *testing.T) {
	srv := newReleaseServer(t, "not-json")
	defer srv.Close()

	assert.Nil(t, checkAt(srv.URL, "v2.0.0"))
}

func TestCheckAt_ReturnsNilWhenServerUnreachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	assert.Nil(t, checkAt(srv.URL, "v2.0.0"))
}

func TestCheckAt_ReturnsNilOnMissingTag(t *testing.T) {
	srv := newReleaseServer(t, `{"assets":[]}`)
	defer srv.Close()

	assert.Nil(t, checkAt(srv.URL, "v2.0.0"))
}

func TestAssetURL_PicksDesktopAssetsNotCLIAssets(t *testing.T) {
	assert.EqualValues(t,
		"https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-ui-darwin-universal.zip",
		assetURL("darwin", mittiV630Assets))

	assert.EqualValues(t,
		"https://github.com/SafetyCulture/mitti-exporter/releases/download/v6.3.0/mitti-exporter-ui-windows-x86_64.exe",
		assetURL("windows", mittiV630Assets))
}

func TestAssetURL_ReturnsEmptyWhenDesktopAssetMissing(t *testing.T) {
	// only CLI assets published: nothing downloadable for the desktop app
	assert.Empty(t, assetURL("darwin", mittiV630Assets[:3]))
	assert.Empty(t, assetURL("windows", mittiV630Assets[:3]))

	// unsupported platform
	assert.Empty(t, assetURL("linux", mittiV630Assets))
}

func TestIsNewer(t *testing.T) {
	assert.True(t, isNewer("v2.0.0", "v6.3.0"))
	assert.False(t, isNewer("v6.3.0", "v6.3.0"))
	assert.False(t, isNewer("v6.4.0", "v6.3.0"))

	// unparseable current version: report newer and defer to ShouldUpdate
	assert.True(t, isNewer("not-a-version", "v6.3.0"))
	// unparseable latest version: nothing to report
	assert.False(t, isNewer("v2.0.0", "not-a-version"))
}
