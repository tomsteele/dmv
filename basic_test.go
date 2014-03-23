package dmv

import (
	"encoding/base64"
	"fmt"
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_BasicAuth(t *testing.T) {
	res := httptest.NewRecorder()
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("gopher:golf"))
	m := martini.Classic()
	m.Get("/protected", AuthBasic(), func(w http.ResponseWriter, req *http.Request, b *Basic) {
		fmt.Fprintf(w, "hi %s %s", b.Username, b.Password)
	})
	r, _ := http.NewRequest("GET", "/protected", nil)
	m.ServeHTTP(res, r)
	if res.Code != 401 {
		t.Error("Response not 401")
	}
	if strings.Contains(res.Body.String(), "hi") {
		t.Error("Auth block failed")
	}
	res = httptest.NewRecorder()
	r.Header.Set("Authorization", auth)
	m.ServeHTTP(res, r)
	if res.Code == 401 {
		t.Error("Response is 401")
	}
	if res.Body.String() != "hi gopher golf" {
		t.Error("Auth failed, got: ", res.Body.String())
	}
}
