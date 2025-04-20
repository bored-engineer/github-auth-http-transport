package ghauth

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/oauth2"
)

// Transport returns an http.RoundTripper that adds GitHub authentication to each HTTP request.
// If base is nil, http.DefaultTransport is used.
// It will try (in order): Environment, App, OAuth, Netrc, CLI
func Transport(ctx context.Context, base http.RoundTripper) (http.RoundTripper, error) {
	if base == nil {
		base = http.DefaultTransport
	}

	token, err1 := Environment()
	if err1 == nil {
		return &oauth2.Transport{
			Base:   base,
			Source: oauth2.StaticTokenSource(token),
		}, nil
	}

	ts, err2 := App(ctx, "", "", "")
	if err2 == nil {
		return &oauth2.Transport{
			Base:   base,
			Source: ts,
		}, nil
	}

	basicauth, err3 := OAuth(base, "", "")
	if err3 == nil {
		return basicauth, nil
	}

	token, err4 := Netrc("", "")
	if err4 == nil {
		return &oauth2.Transport{
			Base:   base,
			Source: oauth2.StaticTokenSource(token),
		}, nil
	}

	token, err5 := CLI("", "")
	if err5 == nil {
		return &oauth2.Transport{
			Base:   base,
			Source: oauth2.StaticTokenSource(token),
		}, nil
	}

	return nil, errors.Join(err1, err2, err3, err4, err5)
}
