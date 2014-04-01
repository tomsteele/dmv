package dmv

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	ghProfileURL = "https://api.github.com/user"
)

// Github stores the access and refresh tokens along with the users profile.
type Github struct {
	Errors       []error
	AccessToken  string
	RefreshToken string
	Profile      GithubProfile
}

// GithubProfile stores information about the user from Github.
type GithubProfile struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Login   string `json:"login"`
	HTMLURL string `json:"html_url"`
	Email   string `json:"email"`
}

// AuthGithub authenticates users using Github and OAuth2.0. After handling
// a callback request, a request is made to get the users Github profile
// and a Github struct will be mapped to the current request context.
//
// This function should be called twice in each application, once on the login
// handler and once on the callback handler.
//
//
//     package main
//
//     import (
//         "github.com/go-martini/martini"
//         "github.com/martini-contrib/sessions"
//         "net/http"
//     )
//
//     func main() {
//         ghOpts := &dmv.OAuth2.0Options{
//             ClientID: "oauth_id",
//             ClientSecret: "oauth_secret",
//             RedirectURL: "http://host:port/auth/callback/github",
//         }
//
//         m := martini.Classic()
//         store := sessions.NewCookieStore([]byte("secret123"))
//         m.Use(sessions.Sessions("my_session", store))
//
//         m.Get("/", func(s sessions.Session) string {
//             return "hi" + s.Get("userID")
//         })
//         m.Get("/auth/github", dmv.AuthGithub(ghOpts))
//         m.Get("/auth/callback/github", dmv.AuthGithub(ghOpts), func(gh *dmv.Github, req *http.Request, w http.ResponseWriter) {
//             // Handle any errors.
//             if len(gh.Errors) > 0 {
//                 http.Error(w, "Oauth failure", http.StatusInternalServerError)
//                 return
//             }
//             // Do something in a database to create or find the user by the Github profile id.
//             user := findOrCreateByGithubID(gh.Profile.ID)
//             s.Set("userID", user.ID)
//             http.Redirect(w, req, "/", http.StatusFound)
//         })
//     }
func AuthGithub(opts *OAuth2Options) martini.Handler {
	opts.AuthURL = "https://github.com/login/oauth/authorize"
	opts.TokenURL = "https://github.com/login/oauth/access_token"
	config := &oauth.Config{
		ClientId:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURL:  opts.RedirectURL,
		Scope:        strings.Join(opts.Scopes, ","),
		AuthURL:      opts.AuthURL,
		TokenURL:     opts.TokenURL,
	}

	transport := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
	}

	cbPath := ""
	if u, err := url.Parse(opts.RedirectURL); err == nil {
		cbPath = u.Path
	}
	return func(r *http.Request, w http.ResponseWriter, c martini.Context) {
		if r.URL.Path != cbPath {
			http.Redirect(w, r, transport.Config.AuthCodeURL(""), http.StatusFound)
			return
		}
		gh := &Github{}
		defer c.Map(gh)
		code := r.FormValue("code")
		tk, err := transport.Exchange(code)
		if err != nil {
			gh.Errors = append(gh.Errors, err)
			return
		}
		gh.AccessToken = tk.AccessToken
		gh.RefreshToken = tk.RefreshToken
		resp, err := transport.Client().Get(ghProfileURL)
		if err != nil {
			gh.Errors = append(gh.Errors, err)
			return
		}
		defer resp.Body.Close()
		profile := &GithubProfile{}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			gh.Errors = append(gh.Errors, err)
			return
		}
		if err := json.Unmarshal(data, profile); err != nil {
			gh.Errors = append(gh.Errors, err)
			return
		}
		gh.Profile = *profile
		return
	}
}
