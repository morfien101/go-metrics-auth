package webengine

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeRedis struct {
	createCredentials func() ([]string, error)
	validateAuth      func(username, password string) (bool, error)
	getEndpoint       func() (string, error)
}

func (fr *fakeRedis) CreateCredentials() ([]string, error) {
	return fr.createCredentials()
}
func (fr *fakeRedis) ValidateAuth(username, password string) (bool, error) {
	return fr.validateAuth(username, password)
}
func (fr *fakeRedis) GetEndpoint() (string, error) {
	return fr.getEndpoint()
}

func TestValidate(t *testing.T) {
	username := "$£!good_user!£$"
	password := "£$%^&*!ANDHESakshdajnd&&*()!!!"

	re := &fakeRedis{
		validateAuth: func(u, p string) (bool, error) {
			if u == username && p == password {
				return true, nil
			}
			return false, nil
		},
	}

	req, _ := http.NewRequest("POST", "http://localhost/auth", nil)
	queries := req.URL.Query()
	queries.Add("username", username)
	queries.Add("password", password)
	req.URL.RawQuery = queries.Encode()

	we := &WebEngine{
		redisEngine: re,
	}
	we.loadRoutes()
	w := httptest.NewRecorder()
	we.ServeHTTP(w, req)

	if w.Result().StatusCode != 200 {
		t.Logf("Status Code is not expected. Got %d, Want: %d", w.Result().StatusCode, 200)
		t.Fail()
	}

	queries = req.URL.Query()
	queries.Del("username")
	queries.Add("username", "Potatoes")
	req.URL.RawQuery = queries.Encode()
	w = httptest.NewRecorder()
	we.ServeHTTP(w, req)

	if w.Result().StatusCode != 401 {
		t.Logf("Status Code is not expected. Got %d, Want: %d", w.Result().StatusCode, 401)
		t.Fail()
	}

	re.validateAuth = func(string, string) (bool, error) { return false, fmt.Errorf("Fake Error") }
	w = httptest.NewRecorder()
	we.ServeHTTP(w, req)

	if w.Result().StatusCode != 500 {
		t.Logf("Status Code is not expected. Got %d, Want: %d", w.Result().StatusCode, 500)
		t.Fail()
	}
}
