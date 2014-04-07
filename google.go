package dmv

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-martini/martini"
)

var (
	googleProfileURL = "https://www.googleapis.com/oauth2/v1/userinfo"
)

// Google stores the access and refresh tokens along with the user profile.
type Google struct {
	Errors       []error
	AccessToken  string
	RefreshToken string
	Profile      GoogleProfile
}

// GoogleProfile stores information from the users google+ profile.
type GoogleProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"name"`
	FamilyName  string `json:"family_name"`
	GivenName   string `json:"given_name"`
	Email       string `json:"email"`
}

// AuthGoogle authenticates users using Google and OAuth2.0. After handling
// a callback request, a request is made to get the users Google profile
// and a Google struct will be mapped to the current request context.
//
// This function should be called twice in each application, once on the login
// handler and once on the callback handler.
//
//     package main
//
//     import (
//         "github.com/go-martini/martini"
//         "github.com/martini-contrib/sessions"
//         "github.com/thomasjsteele/dmv"
//         "net/http"
//     )
//
//     func main() {
//         googleOpts := &dmv.OAuth2Options{
//             ClientID: "oauth_id",
//             ClientSecret: "oauth_secret",
//             RedirectURL: "http://host:port/auth/callback/google",
//             Scopes:      []string{"https://www.googleapis.com/auth/userinfo.email",
//                                   "https://www.googleapis.com/auth/userinfo.profile"},
//         }
//
//         m := martini.Classic()
//         store := sessions.NewCookieStore([]byte("secret123"))
//         m.Use(sessions.Sessions("my_session", store))
//
//         m.Get("/", func(s sessions.Session) string {
//             return "hello" + s.Get("userID")
//         })
//         m.Get("/auth/google", dmv.AuthGoogle(googleOpts))
//         m.Get("/auth/callback/google", dmv.AuthGoogle(googleOpts), func(goog *dmv.Google, req *http.Request, w http.ResponseWriter) {
//             // Handle any errors.
//             if len(goog.Errors) > 0 {
//                 http.Error(w, "OAuth failure", http.StatusInternalServerError)
//                 return
//             }
//             // Do something in a database to create or find the user by the Google profile id.
//             s.Set("userID", goog.Profile.ID)
//             http.Redirect(w, req, "/", http.StatusFound)
//         })
//     }
func AuthGoogle(opts *OAuth2Options) martini.Handler {
	opts.AuthURL = "https://accounts.google.com/o/oauth2/auth"
	opts.TokenURL = "https://accounts.google.com/o/oauth2/token"

	return func(r *http.Request, w http.ResponseWriter, c martini.Context) {
		transport := makeTransport(opts, r)
		cbPath := ""
		if u, err := url.Parse(transport.Config.RedirectURL); err == nil {
			cbPath = u.Path
		}
		if r.URL.Path != cbPath {
			http.Redirect(w, r, transport.Config.AuthCodeURL(""), http.StatusFound)
			return
		}
		goog := &Google{}
		defer c.Map(goog)
		code := r.FormValue("code")
		tk, err := transport.Exchange(code)
		if err != nil {
			goog.Errors = append(goog.Errors, err)
			return
		}
		goog.AccessToken = tk.AccessToken
		goog.RefreshToken = tk.RefreshToken
		resp, err := transport.Client().Get(googleProfileURL)
		if err != nil {
			goog.Errors = append(goog.Errors, err)
			return
		}
		defer resp.Body.Close()
		profile := &GoogleProfile{}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			goog.Errors = append(goog.Errors, err)
			return
		}
		if err := json.Unmarshal(data, profile); err != nil {
			goog.Errors = append(goog.Errors, err)
			return
		}
		goog.Profile = *profile
		return
	}
}
