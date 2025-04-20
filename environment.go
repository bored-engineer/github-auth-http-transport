package ghauth

import (
	"errors"
	"os"

	"golang.org/x/oauth2"
)

// Host returns the host name of the GitHub server by reading the GH_HOST environment variable.
// If the variable is not set, it defaults to "github.com".
func Host() string {
	host := os.Getenv("GH_HOST")
	if host == "" {
		host = "github.com"
	}
	return host
}

// Environment reads a static GitHub token from the environment
// It tries (in order): GH_TOKEN, GITHUB_TOKEN, GH_ENTERPRISE_TOKEN, GITHUB_ENTERPRISE_TOKEN
func Environment() (*oauth2.Token, error) {
	for _, key := range []string{
		"GH_TOKEN",
		"GITHUB_TOKEN",
		"GH_ENTERPRISE_TOKEN",
		"GITHUB_ENTERPRISE_TOKEN",
	} {
		if token := os.Getenv(key); token != "" {
			return &oauth2.Token{AccessToken: token, TokenType: "Bearer"}, nil
		}
	}
	return nil, errors.New("no token found in environment variables")
}
