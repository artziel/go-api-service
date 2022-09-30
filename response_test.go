package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHttpTest(method string, url string, handlerFnc func(w http.ResponseWriter, r *http.Request)) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFnc)

	handler.ServeHTTP(rr, req)

	return rr, nil
}

func TestRespondWithJSONError(t *testing.T) {

	rr, err := newHttpTest("GET", "/TestRespondWithJSONMessage", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("internal server error")
		RespondWithJSONError(w, http.StatusInternalServerError, err)
	}))
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"internal server error"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \n\t got %v\n\twant %v", rr.Body.String(), expected)
	}
}

func TestRespondWithJSONMessage(t *testing.T) {

	rr, err := newHttpTest("GET", "/TestRespondWithJSONMessage", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSONMessage(w, http.StatusOK, "sample json message")
	}))
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"sample json message"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \n\t got %v\n\twant %v", rr.Body.String(), expected)
	}
}

func TestRespondWithJSON(t *testing.T) {

	rr, err := newHttpTest("GET", "/TestRespondWithJSON", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"message": "sample message",
			"id":      1,
			"price":   1.9,
		})
	}))
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":1,"message":"sample message","price":1.9}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \n\t got %v\n\twant %v", rr.Body.String(), expected)
	}
}

func TestRespondWithJSONHMAC(t *testing.T) {
	secretKey := "4234kxzjcjj3@nxnxbcvsjfj"
	payload := map[string]interface{}{
		"message": "sample message",
		"id":      1,
		"price":   1.9,
	}

	rr, err := newHttpTest(
		"GET", "/TestRespondWithJSON",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RespondWithJSONHMAC(w, http.StatusOK, payload, secretKey)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	header := rr.Header()
	resHash := header.Get("Service-Content-Hash")
	if resHash == "" {
		t.Errorf("RespondWithJSONHMAC respond without a \"content-hash\" header")
	}
	newHash := NewHash(rr.Body.String(), secretKey)
	if resHash != newHash {
		t.Errorf("RespondWithJSONHMAC response hash did not match with new hashed body")
	}
}
