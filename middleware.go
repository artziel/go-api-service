package ApiService

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func MiddlewareAccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			log.Printf(
				"PID %d, Rotines %d - [%s] from IP: %s - URL: %s\n",
				os.Getpid(), runtime.NumGoroutine(),
				request.Method, request.RemoteAddr, request.URL,
			)
			next.ServeHTTP(response, request)
		})
}

func MiddlewareRestrictToLocal(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(response http.ResponseWriter, request *http.Request) {
			ip := request.RemoteAddr[:strings.LastIndex(request.RemoteAddr, ":")]
			if ip != "127.0.0.1" {
				RespondWithJSONMessage(response, http.StatusForbidden, "the request endpoint is restricted")
			} else {
				next.ServeHTTP(response, request)
			}
		})
}
