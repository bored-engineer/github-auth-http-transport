package ghauth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bored-engineer/basicauth"
)

// OAuth returns a basicauth.Transport using a Github OAuth App's client ID and client secret.
// If clientID is empty, it will look for the GH_CLIENT_ID environment variable.
// If clientSecret is empty, it will look for the GH_CLIENT_SECRET environment variable.
func OAuth(
	base http.RoundTripper,
	clientID string,
	clientSecret string,
) (*basicauth.Transport, error) {
	if clientID == "" {
		clientID = os.Getenv("GH_CLIENT_ID")
		if clientID == "" {
			return nil, fmt.Errorf("GH_CLIENT_ID is not set")
		}
	}
	if clientSecret == "" {
		clientSecret = os.Getenv("GH_CLIENT_SECRET")
		if clientSecret == "" {
			return nil, fmt.Errorf("GH_CLIENT_SECRET is not set")
		}
	}
	return &basicauth.Transport{
		Username: clientID,
		Password: clientSecret,
		Base:     base,
	}, nil
}
