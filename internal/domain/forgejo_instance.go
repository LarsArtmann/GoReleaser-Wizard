package domain

// Forgejo instance coordinates for the LarsArtmann fleet. The instance is
// LAN-internal by design; public consumers (proxy.golang.org, Homebrew,
// binary downloads) are served by the GitHub push-mirror shadow.
const (
	// ForgejoInstanceURL is the root URL of the primary Forgejo instance.
	ForgejoInstanceURL = "https://forgejo.home.lan"
	// ForgejoAPIURL is the v1 API root used by GoReleaser's gitea_urls.
	ForgejoAPIURL = ForgejoInstanceURL + "/api/v1"
)
