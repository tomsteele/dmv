package dmv

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-martini/martini"
)

func TestLoginRedirect(t *testing.T) {
	recorder := httptest.NewRecorder()
	googleOpts := &OAuth2Options{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		RedirectURL:  "refresh_url",
		Scopes:       []string{"x", "y"},
	}
	m := testMartini()
	m.Get("/auth/google", AuthGoogle(googleOpts))

	r, _ := http.NewRequest("GET", "/auth/google", nil)
	m.ServeHTTP(recorder, r)

	location := recorder.HeaderMap["Location"][0]
	if recorder.Code != 302 {
		t.Errorf("Not being redirected to the auth page.")
	}
	if location != "https://accounts.google.com/o/oauth2/auth?access_type=&approval_prompt=&client_id=client_id&redirect_uri=refresh_url&response_type=code&scope=x+y&state=" {
		t.Errorf("Not being redirected to the right page, %v found", location)
	}
}

func TestLoginRedirectFunc(t *testing.T) {
	recorder := httptest.NewRecorder()
	googleOpts := &OAuth2Options{
		RedirectFunc: RedirectRelativeFunc("/auth/callback/google"),
	}
	m := testMartini()
	m.Get("/auth/google", AuthGoogle(googleOpts))

	r, _ := http.NewRequest("GET", "/auth/google", nil)
	m.ServeHTTP(recorder, r)

	location := recorder.HeaderMap["Location"][0]
	if recorder.Code != 302 {
		t.Errorf("Not being redirected to the auth page.")
	}
	u, err := url.Parse(location)
	if err != nil {
		t.Fatal(err)
	}
	uri, err := url.QueryUnescape(u.Query().Get("redirect_uri"))
	if err != nil {
		t.Fatal(err)
	}
	if uri != "http://localhost/auth/callback/google" {
		t.Errorf("Not being redirected to the right URL, %q found", uri)
	}
}

func testMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	return &martini.ClassicMartini{m, r}
}
