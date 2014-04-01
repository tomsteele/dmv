package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/tomsteele/dmv"
	"net/http"
)

func main() {
	googleOpts := &dmv.OAuth2Options{
		ClientID:     "CLIENTID",
		ClientSecret: "CLIENTSECRET",
		RedirectURL:  "http://localhost:3000/auth/callback/google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}

	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))
	m.Use(render.Renderer())

	m.Get("/", func(s sessions.Session, r render.Render) {
		if uid := s.Get("userID"); uid == nil {
			r.Redirect("/login", 302)
		} else {
			r.HTML(200, "home", uid.(string))
		}
	})

	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})

	m.Get("/auth/google", dmv.AuthGoogle(googleOpts))
	m.Get("/auth/callback/google", dmv.AuthGoogle(googleOpts), func(goog *dmv.Google, r render.Render, s sessions.Session, w http.ResponseWriter) {
		// Handle any errors.
		if len(goog.Errors) > 0 {
			http.Error(w, "Oauth failure", http.StatusInternalServerError)
			return
		}
		// Do something in a database to create or find the user by the Google profile id.
		// for now just pass the google display name
		s.Set("userID", goog.Profile.ID)
		r.HTML(200, "home", goog.Profile.DisplayName)
	})

	m.Run()
}
