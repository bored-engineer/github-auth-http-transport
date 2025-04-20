package ghauth

import (
	"golang.org/x/oauth2"
)

// Token creates a new oauth2.Token with the given token string.
// The token type is set to "Bearer" by default.
func Token(token string) *oauth2.Token {
	return &oauth2.Token{
		TokenType:   "Bearer",
		AccessToken: token,
	}
}

// TokenSource creates a new oauth2.TokenSource using the provided token.
func TokenSource(token *oauth2.Token) oauth2.TokenSource {
	return oauth2.StaticTokenSource(token)
}

// TokenTransport creates a new oauth2.Transport using the provided token.
func TokenTransport(token *oauth2.Token) *oauth2.Transport {
	return &oauth2.Transport{
		Source: TokenSource(token),
	}
}
