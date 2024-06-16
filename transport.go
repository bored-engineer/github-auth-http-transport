package githubauth

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"golang.org/x/oauth2"
)

// Transport returns an http.RoundTripper that makes authenticated requests to the GitHub API.
func Transport(base http.RoundTripper) (http.RoundTripper, error) {
	if base == nil {
		base = http.DefaultTransport
	}
	// The GitHub CLI supports multiple options for tokens via environment variables
	var token string
	for _, key := range []string{"GH_TOKEN", "GITHUB_TOKEN", "GH_ENTERPRISE_TOKEN", "GITHUB_ENTERPRISE_TOKEN"} {
		if value := os.Getenv(key); value != "" {
			token = value
			break
		}
	}
	// If we found a token, we're done!
	if token != "" {
		return &oauth2.Transport{
			Base: base,
			Source: oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: token,
			}),
		}, nil
	}
	// If there is a GITHUB_APP_ID and GITHUB_APP_INSTALLATION_ID, we can use the GitHub App installation token
	appIDStr, installationIDStr := os.Getenv("GITHUB_APP_ID"), os.Getenv("GITHUB_APP_INSTALLATION_ID")
	if appIDStr != "" && installationIDStr != "" {
		// Ensure the environment variables can be parsed
		appID, err := strconv.ParseInt(appIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse GITHUB_APP_ID(%q): %w", appIDStr, err)
		}
		installationID, err := strconv.ParseInt(installationIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse GITHUB_APP_INSTALLATION_ID(%q): %w", installationIDStr, err)
		}
		// Read the private key from a file or environment variable
		var privateKeyBytes []byte
		if path := os.Getenv("GITHUB_APP_PRIVATE_KEY_FILE"); path != "" {
			var err error
			privateKeyBytes, err = os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("could not read GITHUB_APP_PRIVATE_KEY_FILE(%q): %w", path, err)
			}
		} else {
			if value := os.Getenv("GITHUB_APP_PRIVATE_KEY"); value != "" {
				privateKeyBytes = []byte(value)
			} else {
				return nil, fmt.Errorf("GITHUB_APP_PRIVATE_KEY_FILE or GITHUB_APP_PRIVATE_KEY must be set")
			}
		}
		// Create the GitHub App client
		transport, err := ghinstallation.New(base, appID, installationID, privateKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("ghinstallation.New failed: %w", err)
		}
		// Default to github.com but support GitHub Enterprise use-cases as well
		if value := os.Getenv("GH_HOST"); value != "" {
			transport.BaseURL = "https://" + strings.TrimPrefix(value, "https://")
		}
		return transport, nil
	}
	// Final fallback check the GitHub CLI
	if ghCLI, err := exec.LookPath("gh"); err == nil {
		cmd := exec.Command(ghCLI, "auth", "token")
		if value := os.Getenv("GH_HOST"); value != "" {
			cmd.Args = append(cmd.Args, "--hostname", value)
		}
		out, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("%v failed: %w", cmd.Args, err)
		}
		return &oauth2.Transport{
			Base: base,
			Source: oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: strings.TrimSpace(string(out)),
			}),
		}, nil
	}
	return nil, fmt.Errorf("no GitHub token found")
}
