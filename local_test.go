package dmv

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func Test_AuthLocal(t *testing.T) {
	m := martini.Classic()
	m.Post("/login", AuthLocal(&LocalOptions{}), func(l *Local, req *http.Request, w http.ResponseWriter) {
		if len(l.Errors) > 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%s %s", l.Username, l.Password)
		return
	})
	user := "gophers"
	pass := "rule"
	data := url.Values{}
	data.Set("username", user)
	data.Set("password", pass)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
	m.ServeHTTP(res, req)
	if res.Code == 400 {
		t.Error("AuthLocal failed to parse valid username and password")
	}
	if res.Body.String() != user+" "+pass {
		t.Error("AuthLocal did not return the correct username and password")
	}
}
