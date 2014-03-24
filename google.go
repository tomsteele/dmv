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
  // google oauth2 endpoint
  googleProfileURL = "https://accounts.google.com/o/oauth2/auth"
)

// Google holds the access and refresh tokens along with the user profile
type Google struct {
  Errors       []error
  AccessToken  string
  RefreshToken string
  Profile      GoogleProfile
}
