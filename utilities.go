package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

/*
Encode a golang error in JSON format and write to the response writer
*/
func RespondWithJSONError(w http.ResponseWriter, code int, err error) {
	RespondWithJSON(w, code, map[string]string{"message": err.Error()})
}

/*
Encode a text message in JSON format and write to the response writer
*/
func RespondWithJSONMessage(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"message": message})
}

/*
Encode a interface data in JSON format and write to the response writer
*/
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

/*
Reads the form information of a request and assigns the form data to its corresponding structure fields
*/
func FormToStruct(r *http.Request, model interface{}) error {
	v := reflect.ValueOf(model)

	if v.Kind() != reflect.Ptr {
		return errors.New("function \"FormToStruct\" expect an struct ptr")
	}

	for i := 0; i < v.Elem().NumField(); i++ {
		tag := v.Elem().Type().Field(i).Tag.Get("form")
		if tag != "" {
			value := r.FormValue(tag)
			if value != "" {
				switch v.Elem().Type().Field(i).Type.Kind() {
				case reflect.Bool:
					val, _ := strconv.ParseBool(value)
					v.Elem().Field(i).SetBool(val)
				case reflect.String:
					v.Elem().Field(i).SetString(value)
				case reflect.Float32, reflect.Float64:
					val, _ := strconv.ParseFloat(value, 64)
					v.Elem().Field(i).SetFloat(val)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					val, _ := strconv.Atoi(value)
					v.Elem().Field(i).SetInt(int64(val))
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					val, _ := strconv.Atoi(value)
					v.Elem().Field(i).SetUint(uint64(val))
				}
			}
		}
	}

	return nil
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

func xffIP(r *http.Request) string {
	var remoteIP string
	var xff string = strings.Trim(r.Header.Get("X-Forwarded-For"), ",")

	if len(xff) != 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
	}
	return remoteIP
}

func xriIP(r *http.Request) string {
	var remoteIP string
	var xri string = r.Header.Get("X-Real-Ip")

	if ip := net.ParseIP(xri); ip != nil {
		remoteIP = ip.String()
	}

	return remoteIP
}

func remoteAddr(r *http.Request) string {
	ip := ""
	var parts []string = strings.Split(r.RemoteAddr, ":")

	if len(parts) == 2 {
		ip = parts[0]
	}

	return ip
}

func GetRealIPAddr(r *http.Request) string {
	var remoteIP string = remoteAddr(r)
	var xff string = xffIP(r)
	var xri string = xriIP(r)

	if xff != "" {
		remoteIP = xff
	}
	if xri != "" {
		remoteIP = xff
	}

	return remoteIP
}
