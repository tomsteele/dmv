package dmv

import (
	"net/http"
	"strings"

	"github.com/tomsteele/dmv/oauth"
)

// OAuth2Options contains options for complete OAuth2.0.
type OAuth2Options struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	// Accepts a func to generate the redirect URL based on the request. Useful
	// if you want to redirect to a relative path. Takes precedence over
	// RedirectURL if both are set.
	RedirectFunc func(*http.Request) string
	Scopes       []string
	AuthURL      string
	TokenURL     string
}

func RedirectRelativeFunc(path string) func(*http.Request) string {
	return func(req *http.Request) string {
		proto := "http"
		if strings.EqualFold(req.URL.Scheme, "https") ||
			req.TLS != nil ||
			req.Header.Get("X-Forwarded-Proto") == "https" ||
			req.Header.Get("X-SSL-Request") == "on" {
			proto = "https"
		}
		host := req.Host
		if host == "" {
			host = "localhost"
		}
		return proto + "://" + host + path
	}
}

func makeTransport(opts *OAuth2Options, req *http.Request) (transport *oauth.Transport) {
	config := &oauth.Config{
		ClientId:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURL:  opts.RedirectURL,
		Scope:        strings.Join(opts.Scopes, " "),
		AuthURL:      opts.AuthURL,
		TokenURL:     opts.TokenURL,
	}
	if opts.RedirectFunc != nil {
		config.RedirectURL = opts.RedirectFunc(req)
	}

	return &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
	}
}
