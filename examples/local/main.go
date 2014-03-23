// This is a simple example of how the local auth
// can be used. Requires MongoDB running on localhost.
//
// Creates a test user of gopher@gophermail.com
// with the password of Password1
package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/tomsteele/dmv"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

// User of the application.
type User struct {
	ID           bson.ObjectId `bson:"_id"`
	Email        string        `bson:"email"`
	PasswordHash string        `bson:"password_hash"`
}

func main() {
	m := martini.Classic()
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))
	m.Use(render.Renderer())
	m.Use(DB())

	m.Get("/", func(s sessions.Session, r render.Render, db *mgo.Database) {
		if s.Get("userID") == nil {
			r.Redirect("/login", 302)
			return
		}
		// Attempt to find the user by the ID provided by the session.
		u := &User{}
		if err := db.C("users").Find(bson.M{"_id": bson.ObjectIdHex(s.Get("userID").(string))}).One(&u); err != nil {
			// User wasn't found.
			s.Clear()
			r.Redirect("/login", 302)
			return
		}
		r.HTML(200, "home", u.Email)
	})

	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})

	m.Post("/login", dmv.AuthLocal(&dmv.LocalOptions{}), func(s sessions.Session, l *dmv.Local, r render.Render, db *mgo.Database) {
		// There were errors in the request.
		if len(l.Errors) > 0 {
			r.HTML(200, "login", "Invalid username or password!")
			return
		}
		// Attempt to find the user by l.Username.
		u := &User{}
		if err := db.C("users").Find(bson.M{"email": l.Username}).One(&u); err != nil {
			// User wasn't found.
			r.HTML(200, "login", "Invalid username or password!")
			return
		}
		// Compare the password.
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(l.Password)); err != nil {
			// Password was wrong.
			r.HTML(200, "login", "Invalid username or password!")
			return
		}
		// Password was correct. Set the session variable and redirect.
		s.Set("userID", u.ID.Hex())
		r.Redirect("/", 302)
	})

	m.Run()
}

// DB maps a MongoDB session to a request.
func DB() martini.Handler {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	// Remove all existing users and create a test user.
	if _, err := session.DB("dmv").C("users").RemoveAll(bson.M{}); err != nil {
		panic(err)
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("Password1"), 10)
	u := &User{ID: bson.NewObjectId(), Email: "gopher@gophermail.com", PasswordHash: string(hash)}
	if err := session.DB("dmv").C("users").Insert(u); err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB("dmv"))
		defer s.Close()
		c.Next()
	}
}
