package ghauth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jdx/go-netrc"
	"golang.org/x/oauth2"
)

// Netrc attempts to extract a static API token from a .netrc file.
// If path is empty, ~/.netrc is used
// If host is empty, the default host is used
func Netrc(path string, host string) (*oauth2.Token, error) {
	if host == "" {
		host = Host()
	}
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("os.UserHomeDir failed: %w", err)
		}
		path = filepath.Join(home, ".netrc")
	}
	creds, err := netrc.Parse(path)
	if err != nil {
		return nil, err
	}
	if machine := creds.Machine(host); machine != nil {
		if password := machine.Get("password"); password != "" {
			return Token(password), nil
		}
	}
	if machine := creds.Machine("api." + host); machine != nil {
		if password := machine.Get("password"); password != "" {
			return Token(password), nil
		}
	}
	return nil, fmt.Errorf("no token found in %s for host %s", path, host)
}
