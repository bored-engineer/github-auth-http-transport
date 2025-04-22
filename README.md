# GitHub Authentication HTTP Transport [![Go Reference](https://pkg.go.dev/badge/github.com/bored-engineer/github-auth-http-transport.svg)](https://pkg.go.dev/github.com/bored-engineer/github-auth-http-transport)
A Golang [http.RoundTripper](https://pkg.go.dev/net/http#RoundTripper) for injecting [GitHub authentication headers](https://docs.github.com/en/rest/authentication/authenticating-to-the-rest-api?apiVersion=2022-11-28) to GitHub's REST API.

## Usage

Here's an example of using the `Transport` method with [go-github](https://github.com/google/go-github):
```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	ghauth "github.com/bored-engineer/github-auth-http-transport"
	"github.com/google/go-github/v71/github"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	transport, err := ghauth.Transport(ctx, nil)
	if err != nil {
		log.Fatalf("ghauth.Transport failed: %v", err)
	}
	client := github.NewClient(&http.Client{
		Transport: transport,
	})

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		log.Fatalf("(*github.UsersService).Get failed: %v", err)
	}

	fmt.Printf("authenticated as: %s\n", user.GetLogin())
}

```

The `Transport` method will try the following sources (in order):
* [Environment](https://pkg.go.dev/github.com/bored-engineer/github-auth-http-transport#Environment): A static GitHub API token (typically a PAT) found in the current environment variables:
	* `GH_TOKEN`
	* `GITHUB_TOKEN`
	* `GH_ENTERPRISE_TOKEN`
	* `GITHUB_ENTERPRISE_TOKEN`
* [App](https://pkg.go.dev/github.com/bored-engineer/github-auth-http-transport#App): Authentication as a GitHub App configured via the `GH_APP_ID`, `GH_APP_INSTALLATION_ID`, and `GH_APP_PRIVATE_KEY` environment variables. 
* [Basic](https://pkg.go.dev/github.com/bored-engineer/github-auth-http-transport#Basic): HTTP Basic Authentication as configured via the `GH_CLIENT_ID` and `GH_CLIENT_SECRET` environment variables.
* [Netrc](https://pkg.go.dev/github.com/bored-engineer/github-auth-http-transport#Netrc): A `machine` entry for `github.com` or `api.github.com` (or `GH_HOST` if populated) in the `~/.netrc` file
* [CLI](https://pkg.go.dev/github.com/bored-engineer/github-auth-http-transport#CLI): A token obtained by executing `gh auth token`
