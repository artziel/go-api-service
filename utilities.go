package ApiService

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func RespondWithJSONError(w http.ResponseWriter, code int, err error) {
	RespondWithJSON(w, code, map[string]string{"message": err.Error()})
}

func RespondWithJSONMessage(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"message": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func FixFileName(name string) string {
	r := regexp.MustCompile("[^aA-zZ0-9ñÑáÁéÉíÍóÓúÚ._()]+")

	return r.ReplaceAllString(name, "-")
}

func SaveFileFromRequest(r *http.Request, formInputName string, dest string) error {
	file, _, err := r.FormFile(formInputName)
	if err != nil {
		return err
	}
	defer file.Close()

	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, _ = io.Copy(f, file)

	return nil
}

func SaveTmpFileFromRequest(r *http.Request, formInputName string, destFolder string) (string, error) {
	input, handler, err := r.FormFile(formInputName)
	if err != nil {
		return "", err
	}
	defer input.Close()

	file, err := ioutil.TempFile(destFolder, "*."+FixFileName(handler.Filename))
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, _ = io.Copy(file, input)

	return file.Name(), nil
}

func ParseAuthorizationHeader(r *http.Request) string {
	value := r.Header.Get("Authorization")

	if len(value) > 6 {
		prefix := strings.ToLower(strings.TrimSpace(value[:7]))

		if prefix == "bearer" {
			value = strings.TrimSpace(value[7:])
		} else {
			value = r.Header.Get("Authorization")
		}
	}
	return value
}
