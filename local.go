/*Package dmv simple authentication schemes for Martini*/
package dmv

import (
	"errors"
	"github.com/codegangsta/martini"
	"net/http"
)

// Local is mapped to the martini.Context from the martini.Handler
// returned from AuthLocal.
type Local struct {
	Errors   []error
	Username string
	Password string
}

// LocalOptions are used to pass conditional arguments to AuthLocal.
type LocalOptions struct {
	// The form field to represent a username.
	UsernameField string
	// The form field to represent a password.
	PasswordField string
}

// AuthLocal attempts to get a username and password from a request.
//
//     m.Post("/login", dmv.AuthLocal(), func(l *dmv.Local) {
//         if len(l.Errors) > 0 {
//             // Return invalid username or password or perhaps 401.
//         }
//         // Lookup the user by l.Username
//         // Compare password of found user to l.Password
//     })
func AuthLocal(opts *LocalOptions) martini.Handler {
	if opts.UsernameField == "" {
		opts.UsernameField = "username"
	}
	if opts.PasswordField == "" {
		opts.PasswordField = "password"
	}
	return func(req *http.Request, c martini.Context) {
		l := &Local{}
		l.Username = req.FormValue(opts.UsernameField)
		if l.Username == "" {
			l.Errors = append(l.Errors, errors.New("username field not found or empty"))
		}
		l.Password = req.FormValue(opts.PasswordField)
		if l.Password == "" {
			l.Errors = append(l.Errors, errors.New("password field not found or empty"))
		}
		c.Map(l)
	}
}
