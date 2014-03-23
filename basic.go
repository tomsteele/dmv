package dmv

import (
	"encoding/base64"
	"github.com/codegangsta/martini"
	"net/http"
	"strings"
)

// Basic holds a username and password
// from the Authorization header.
type Basic struct {
	Username string
	Password string
}

// AuthBasic attempts to get a username and password
// from an Authorization header. Basic is mapped to the current
// request context. BasicFail will be called if there are errors, the header is empty, or
// a username or password is empty.
//
//
//    m.Get("/protected", AuthBasic(), func(b *dmv.Basic, w http.ResponseWriter) {
//        // Lookup user by b.Username
//        // Compare password to b.Password
//        // If not valid call dmv.FailBasic(w)
//    })
func AuthBasic() martini.Handler {
	return func(req *http.Request, w http.ResponseWriter, c martini.Context) {
		b := &Basic{}
		auth := req.Header.Get("Authorization")
		if auth == "" {
			FailBasic(w)
			return
		}
		data, err := base64.StdEncoding.DecodeString(strings.Replace(auth, "Basic ", "", 1))
		if err != nil {
			FailBasic(w)
			return
		}
		parts := strings.Split(strings.Replace(string(data), "Basic ", "", 1), ":")
		if len(parts) < 2 {
			FailBasic(w)
			return
		}
		b.Username = parts[0]
		b.Password = parts[1]
		c.Map(b)
	}
}

// FailBasic writes the required response headers to prompt
// for basic authentication.
func FailBasic(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
	http.Error(w, "Not Authorized", http.StatusUnauthorized)
}
