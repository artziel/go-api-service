package rest

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestMiddlewareAccessLog(t *testing.T) {

	var str bytes.Buffer

	log.SetOutput(&str)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/MiddlewareAccessLog", nil)

	res := httptest.NewRecorder()

	handler(res, req)

	middle := MiddlewareAccessLog(handler)

	middle.ServeHTTP(res, req)

	result := strings.TrimSuffix(str.String(), "\n")
	match, _ := regexp.MatchString(
		`^[0-9]{4}/[0-9]{2}/[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} PID [0-9]+, Routines [0-9]{1,6} - \[(GET|POST|PUT|DELETE|OPTION)\] from IP: [0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}:[0-9]{1,4} - URL: /\w+$`,
		result,
	)
	if !match {
		t.Errorf("log line do not match regular expresion: \n\t got %v", result)
	}

}

func TestMiddlewareRestrictToLocal(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSONMessage(w, http.StatusOK, "OK")
	})

	req := httptest.NewRequest(http.MethodGet, "/MiddlewareRestrictToLocal", nil)
	req.RemoteAddr = "127.0.0.1:12345"

	resOk := httptest.NewRecorder()

	middle := MiddlewareRestrictToLocal(handler)

	middle.ServeHTTP(resOk, req)

	if status := resOk.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expected := `{"message":"OK"}`
	if resOk.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \n\t got %v\n\twant %v", resOk.Body.String(), expected)
	}

	req.RemoteAddr = "192.168.0.1:12345"
	resFail := httptest.NewRecorder()

	middle.ServeHTTP(resFail, req)

	if status := resFail.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
	expected = `{"message":"the request endpoint is restricted"}`
	if resFail.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \n\t got %v\n\twant %v", resFail.Body.String(), expected)
	}
}
