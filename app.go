package ghauth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"
	"strings"

	oauth2githubapp "github.com/int128/oauth2-github-app"
	"golang.org/x/oauth2"
)

// App returns a oauth2.TokenSource for a given GitHub App installation.
// If appID is empty, it will look for the GH_APP_ID environment variable.
// If installationID is empty, it will look for the GH_APP_INSTALLATION_ID environment variable.
// If privateKeyPath is empty, it will look for the GH_APP_PRIVATE_KEY environment variable.
func App(
	ctx context.Context,
	appID string,
	installationID string,
	privateKeyPath string,
) (oauth2.TokenSource, error) {
	if appID == "" {
		appID = os.Getenv("GH_APP_ID")
		if appID == "" {
			return nil, fmt.Errorf("GH_APP_ID is not set")
		}
	}
	if installationID == "" {
		installationID = os.Getenv("GH_APP_INSTALLATION_ID")
		if installationID == "" {
			return nil, fmt.Errorf("GH_APP_INSTALLATION_ID is not set")
		}
	}
	if privateKeyPath == "" {
		privateKeyPath = os.Getenv("GH_APP_PRIVATE_KEY")
		if privateKeyPath == "" {
			return nil, fmt.Errorf("GH_APP_PRIVATE_KEY is not set")
		}
	}
	var privateKey *rsa.PrivateKey
	if strings.Contains(privateKeyPath, "-BEGIN RSA PRIVATE KEY-") {
		// If privateKeyPath contains the private key directly, parse it.
		var err error
		privateKey, err = oauth2githubapp.ParsePrivateKey([]byte(privateKeyPath))
		if err != nil {
			return nil, fmt.Errorf("oauth2githubapp.ParsePrivateKey failed: %w", err)
		}
	} else {
		// Otherwise, read the private key from the file.
		var err error
		privateKey, err = oauth2githubapp.LoadPrivateKey(privateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("oauth2githubapp.LoadPrivateKey failed: %w", err)
		}
	}
	cfg := oauth2githubapp.Config{
		PrivateKey:     privateKey,
		AppID:          appID,
		InstallationID: installationID,
	}
	if host := Host(); host != "github.com" {
		cfg.BaseURL = "https://" + host + "/api/v3/"
	}
	return cfg.TokenSource(ctx), nil
}
