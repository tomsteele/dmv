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
  // google oauth2 endpoints
  googleProfileURL = "https://www.googleapis.com/oauth2/v1/userinfo"
)

// Google holds the access and refresh tokens along with the user profile
type Google struct {
  Errors       []error
  AccessToken  string
  RefreshToken string
  Profile      GoogleProfile
}

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
// Options:
//  - ClientID        your Google application's client id
//  - ClientSecret    your Google application's client secret
//  - RedirectURL     URL to which Google will redirect the user after granting authorization
//  - Scopes          typical scopes will be: userinfo.email and userinfo.profile
//
//     package main
//
//     import (
//         "github.com/codegangsta/martini"
//         "github.com/martini-contrib/sessions"
//         "github.com/thomasjsteele/dmv"
//         "net/http"
//     )
//
//     func main() {
//         googleOpts := &dmv.OAuth2Options{
//             ClientID: "<CLIENT_ID_HERE",
//             ClientSecret: "<CLIENT_SECRET_HERE",
//             RedirectURL: "http://host:port/auth/callback/google",
//             Scopes:      []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
//         }
//
//         m := martini.Classic()
//         store := sessions.NewCookieStore([]byte("secret123"))
//         m.Use(sessions.Sessions("my_session", store))
//
//         m.Get("/", func(s sessions.Session) string {
//             return "hello visitor"
//         })
//         m.Get("/auth/google", dmv.AuthGoogle(googleOpts))
//         m.Get("/auth/callback/google", dmv.AuthGoogle(googleOpts), func(goog *dmv.Google, req *http.Request, w http.ResponseWriter) {
//             // Handle any errors.
//             if len(goog.Errors) > 0 {
//                 http.Error(w, "Oauth failure", http.StatusInternalServerError)
//                 return
//             }
//             // Do something in a database to create or find the user by the Google profile id.
//             //user := findOrCreateByGoogleID(google.Profile.ID)
//             fmt.Printf("Found Google Profile: %s\n", goog.Profile.ID)
//             http.Redirect(w, req, "/", http.StatusFound)
//         })
//     }

func AuthGoogle(opts *OAuth2Options) martini.Handler {
  opts.AuthURL = "https://accounts.google.com/o/oauth2/auth"
  opts.TokenURL = "https://accounts.google.com/o/oauth2/token"
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
