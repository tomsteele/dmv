package dmv

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"github.com/codegangsta/martini"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	fbProfileURL = "https://graph.facebook.com/me"
)

// Facebook holds the access and refresh tokens along with the users
// profile.
type Facebook struct {
	Errors       []error
	AccessToken  string
	RefreshToken string
	Profile      FacebookProfile
}

// FacebookProfile contains information about the user from facebook.
type FacebookProfile struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	Gender     string `json:"gender"`
	Link       string `json:"link"`
	Email      string `json:"email"`
}

// AuthFacebook authenticates users using Facebook and OAuth2.0. After
// handling a callback request, a request is made to get the users
// facebook profile and a Facebook struct will be mapped to the
// current request context.
//
// This function should be called twice in each application, once
// on the login handler, and once on the callback handler.
//
// Example usage:
//
// package main
//
// import (
//     "github.com/codegangsta/martini"
//     "github.com/martini-contrib/sessions"
//     "net/http"
// )
//
// func main() {
//     fbOpts := &dmv.OAuth2.0Options{
//         ClientID: "oauth_id",
//         ClientSecret: "oauth_secret",
//         RedirectURL: "http://host:port/auth/facebook/callback",
//     }
//
//     m := martini.Classic()
//     store := sessions.NewCookieStore([]byte("secret123"))
//     m.Use(sessions.Sessions("my_session", store))
//
//     m.Get("/", func(s sessions.Session) string {
//         return "hi" + s.ID
//     })
//     m.Get("/auth/facebook", dmv.AuthFacebook(fbOpts))
//     m.Get("/auth/callback/facebook", dmv.AuthFacebook(fbOpts), func(fb *dmv.Facebook, req *http.Request, w http.ResponseWriter) {
//         // Handle any errors.
//         if len(fb.Errors) > 0 {
//             http.Error(w, "Oauth failure", http.StatusInternalServerError)
//             return
//         }
//         // Do something in a database to create or find the user by the facebook profile id.
//         user := findOrCreateByFacebookID(fb.Profile.ID)
//         s.Set("userID", user.ID)
//         http.Redirect(w, req, "/", http.StatusFound)
//     })
// }
//
func AuthFacebook(opts *OAuth2Options) martini.Handler {
	opts.AuthURL = "https://www.facebook.com/dialog/oauth"
	opts.TokenURL = "https://graph.facebook.com/oauth/access_token"
	config := &oauth.Config{
		ClientId:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURL:  opts.RedirectURL,
		Scope:        strings.Join(opts.Scopes, " "),
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
		fb := &Facebook{}
		defer c.Map(fb)
		code := r.FormValue("code")
		tk, err := transport.Exchange(code)
		if err != nil {
			fb.Errors = append(fb.Errors, err)
			return
		}
		fb.AccessToken = tk.AccessToken
		fb.RefreshToken = tk.RefreshToken
		resp, err := transport.Client().Get(fbProfileURL)
		if err != nil {
			fb.Errors = append(fb.Errors, err)
			return
		}
		defer resp.Body.Close()
		profile := &FacebookProfile{}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fb.Errors = append(fb.Errors, err)
			return
		}
		if err := json.Unmarshal(data, profile); err != nil {
			fb.Errors = append(fb.Errors, err)
			return
		}
		fb.Profile = *profile
		return
	}
}
