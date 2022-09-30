package rest

import (
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

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
	req.RemoteAddr = "192.168.0.1:888"
	req.Header.Add("X-Forwarded-For", "192.168.0.1")
	req.Header.Add("X-Real-Ip", "192.168.0.2")

	if ip := GetRealIPAddr(req); ip != "192.168.0.1" {
		t.Errorf("Unexpected IP value: \n\t got %v\n\twant %v", ip, "192.168.0.1")
	}
}

func TestFormToStruct(t *testing.T) {

	type Test struct {
		Int    int     `json:"int" form:"int"`
		Uint   uint    `json:"uint" form:"uint"`
		Float  float64 `json:"float" form:"float"`
		String string  `json:"string" form:"string"`
		Bool   bool    `json:"bool" form:"bool"`
	}

	params := url.Values{}
	params.Add("int", "100")
	params.Add("uint", "1")
	params.Add("float", "1.23")
	params.Add("bool", "true")
	params.Add("string", "this is a test")

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(params.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	r.ParseForm()

	test := Test{}

	if err := FormToStruct(r, Test{}); err == nil {
		t.Errorf("Expected ErrFormToStructPtrExpected error")
	}

	err := FormToStruct(r, &test)
	if err != nil {
		t.Errorf("FormToStruct fail! error: %s", err)
	} else {
		if test.Int != 100 {
			t.Errorf("FormToStruct Unexpected value for field \"Int\":\nGot  %v\nWant 100", test.Int)
		}
		if test.Uint != 1 {
			t.Errorf("FormToStruct Unexpected value for field \"Uint\":\nGot  %v\nWant 1", test.Uint)
		}
		if test.Float != 1.23 {
			t.Errorf("FormToStruct Unexpected value for field \"Float\":\nGot  %v\nWant 1.23", test.Float)
		}
		if !test.Bool {
			t.Errorf("FormToStruct Unexpected value for field \"Bool\":\nGot  %v\nWant true", test.Bool)
		}
		if test.String != "this is a test" {
			t.Errorf("FormToStruct Unexpected value for field \"String\":\nGot  \"%v\"\nWant \"this is a test\"", test.String)
		}
	}

}
