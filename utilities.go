package ApiService

import (
	"encoding/json"
	"io"
	"net"
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
	var err error
	file, _, err := r.FormFile(formInputName)
	if err == nil {
		f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
		if err == nil {
			_, _ = io.Copy(f, file)
		}
		defer f.Close()
	}
	defer file.Close()

	return err
}

func SaveTmpFileFromRequest(r *http.Request, formInputName string, destFolder string) (string, error) {
	var result string
	var err error
	input, handler, err := r.FormFile(formInputName)
	if err == nil {
		file, err := os.CreateTemp(destFolder, "*."+FixFileName(handler.Filename))
		if err == nil {
			_, _ = io.Copy(file, input)
			result = file.Name()
		}
		defer file.Close()
	}
	defer input.Close()

	return result, err
}

func ParseAuthorizationHeader(r *http.Request) string {
	value := r.Header.Get("Authorization")

	if len(value) > 6 {
		prefix := strings.ToLower(strings.TrimSpace(value[:7]))

		if prefix == "bearer" {
			value = strings.TrimSpace(value[7:])
		}
	}
	return value
}

func GetRealIPAddr(r *http.Request) string {

	remoteIP := ""
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}

	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}

	return remoteIP
}
