package dmv

// OAuth2Options contains options for complete OAuth2.0.
type OAuth2Options struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	AuthURL      string
	TokenURL     string
}
