package ghauth

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/oauth2"
)

// CLI retrieves an *oauth2.Token by running the 'gh auth token' command to fetch the token used by the GitHub CLI.
// If path is empty, it will try the GH_PATH environment variable, then to look for the 'gh' command in the system PATH.
// If host is empty, it will use the default GitHub host (github.com).
func CLI(path string, host string) (*oauth2.Token, error) {
	if host == "" {
		host = Host()
	}
	if path == "" {
		path = os.Getenv("GH_PATH")
	}
	if path == "" {
		var err error
		path, err = exec.LookPath("gh")
		if err != nil {
			return nil, fmt.Errorf("gh CLI not found in PATH: %w", err)
		}
	}
	cmd := exec.Command(path, "auth", "token", "--hostname", host)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute %v: %w", cmd, err)
	}
	token := string(bytes.TrimSpace(out))
	if token == "" {
		return nil, fmt.Errorf("%v returned empty token", cmd)
	}
	return Token(token), nil
}
