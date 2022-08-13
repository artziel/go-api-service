package ApiService

import (
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestFixFileName(t *testing.T) {

	fName := "Esta Es Una con (ñ,Ñ,á,Á,â+|@/) Prueba.jpg"
	fixed := FixFileName(fName)
	expected := "Esta-Es-Una-con-(ñ-Ñ-á-Á-)-Prueba.jpg"

	if expected != fixed {
		t.Errorf("function return unexpected value: \n\t got %v\n\twant %v", fixed, expected)
	}
}

func TestSaveFileFromRequest(t *testing.T) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	outputFile := "./test.csv"

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("file", "sample.csv")
		if err != nil {
			t.Error(err)
		}

		w := csv.NewWriter(part)

		w.WriteAll([][]string{
			{"ID", "NAME", "EMAIL", "CELLPHONE"},
		})
	}()

	request := httptest.NewRequest("POST", "/SaveFileFromRequest", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := SaveFileFromRequest(r, "file", outputFile)
		if err != nil {
			RespondWithJSONError(w, http.StatusInternalServerError, err)
		} else {
			RespondWithJSONMessage(w, http.StatusOK, "upload success")
		}
	})
	handler.ServeHTTP(rr, request)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"upload success"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \n\t got %v\n\twant %v", rr.Body.String(), expected)
	}

	if _, err := os.Stat(outputFile); errors.Is(err, os.ErrNotExist) {
		t.Errorf("uploaded file \"%v\" not saved", outputFile)
	} else {

		csvFile, err := os.Open(outputFile)
		if err != nil {
			t.Errorf("uploaded csv file \"%v\" can not be open", outputFile)
		} else {
			defer csvFile.Close()
			if csvLines, err := csv.NewReader(csvFile).ReadAll(); err != nil {
				t.Errorf("uploaded csv file \"%v\" can not be read", outputFile)
			} else {
				if csvLines[0][0] != "ID" || csvLines[0][1] != "NAME" || csvLines[0][2] != "EMAIL" || csvLines[0][3] != "CELLPHONE" {
					t.Errorf("unexpected uploaded csv file \"%v\" content", outputFile)
				}
			}
		}

		os.Remove(outputFile)
	}
}

func TestSaveTmpFileFromRequest(t *testing.T) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)
	outputFile := ""

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("file", "sample.csv")
		if err != nil {
			t.Error(err)
		}

		w := csv.NewWriter(part)

		w.WriteAll([][]string{
			{"ID", "NAME", "EMAIL", "CELLPHONE"},
		})
	}()

	request := httptest.NewRequest("POST", "/SaveFileFromRequest", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName, err := SaveTmpFileFromRequest(r, "file", "./")
		if err != nil {
			RespondWithJSONError(w, http.StatusInternalServerError, err)
		} else {
			RespondWithJSONMessage(w, http.StatusOK, "upload success")
			outputFile = fileName
		}
	})
	handler.ServeHTTP(rr, request)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message":"upload success"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: \n\t got %v\n\twant %v", rr.Body.String(), expected)
	}

	if _, err := os.Stat(outputFile); errors.Is(err, os.ErrNotExist) {
		t.Errorf("uploaded file \"%v\" not saved", outputFile)
	} else {

		csvFile, err := os.Open(outputFile)
		if err != nil {
			t.Errorf("uploaded csv file \"%v\" can not be open", outputFile)
		} else {
			defer csvFile.Close()
			if csvLines, err := csv.NewReader(csvFile).ReadAll(); err != nil {
				t.Errorf("uploaded csv file \"%v\" can not be read", outputFile)
			} else {
				if csvLines[0][0] != "ID" || csvLines[0][1] != "NAME" || csvLines[0][2] != "EMAIL" || csvLines[0][3] != "CELLPHONE" {
					t.Errorf("unexpected uploaded csv file \"%v\" content", outputFile)
				}
			}
		}

		os.Remove(outputFile)
	}
}

func TestParseAuthorizationHeader(t *testing.T) {

	req, err := http.NewRequest("GET", "/ParseAuthorizationHeader", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer 12345")

	if token := ParseAuthorizationHeader(req); token != "12345" {
		t.Errorf("Unexpected Authorization header value: \n\t got %v\n\twant %v", token, "12345")
	}

	req.Header.Set("Authorization", "12345")

	if token := ParseAuthorizationHeader(req); token != "12345" {
		t.Errorf("Unexpected Authorization header value: \n\t got %v\n\twant %v", token, "12345")
	}

	req.Header.Set("Authorization", "")

	if token := ParseAuthorizationHeader(req); token != "" {
		t.Errorf("Unexpected Authorization header value: \n\t got %v\n\twant %v", token, "")
	}

	req.Header.Set("Authorization", "Bearer")

	if token := ParseAuthorizationHeader(req); token != "Bearer" {
		t.Errorf("Unexpected Authorization header value: \n\t got %v\n\twant %v", token, "")
	}
}

func TestGetRealIPAddr(t *testing.T) {

	req, err := http.NewRequest("GET", "/ParseAuthorizationHeader", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("X-Forwarded-For", "192.168.0.1")
	req.Header.Add("X-Real-Ip", "192.168.0.2")

	if ip := GetRealIPAddr(req); ip != "192.168.0.1" {
		t.Errorf("Unexpected IP value: \n\t got %v\n\twant %v", ip, "192.168.0.1")
	}

}
