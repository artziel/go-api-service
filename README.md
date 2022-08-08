[![N|Solid](https://www.alus.com.mx/assets/images/logo.svg)](https://www.alus.com.mx/)
# Golang API Service
Golang version 1.18

Artziel Narvaiza <artziel@alus.com.mx>

This is a private repo, you must setup your enviroment to be able to use this repo. Search Google for instructions on how to use Private Repos with Golang.

### Features
- Simplify API creation for Restfull services 
- Easy to understand configuration
- Include JWT for authentication

### Dependencies
- github.com/golang-jwt/jwt
- github.com/gorilla/context
- github.com/gorilla/mux
- gitlab.com/alus/go-security
- golang.org/x/crypto

Get the package
```bash
go get gitlab.com/alus/go-api-service
```

Use example:
```golang
package main

import (
	"fmt"
	"net/http"
	"time"

	ApiService "gitlab.com/alus/go-api-service"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	ApiService.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"Message": "Welcome to API Service sample in Golang!",
	})
}

func sleep10seconds(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	ApiService.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"Message": "Wait for 10 seconds!",
	})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Not Found Error\n")
}

func main() {
	srv := ApiService.NewService(
		"Sample Restful Service",
		"1.0.0b",
		ApiService.ServiceConfig{
			Port: 8080,
		},
	)

	r := srv.Router()
	r.HandleFunc("/welcome", welcome).Methods("GET")
	r.HandleFunc("/sleep-10-seconds", sleep10seconds).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(NotFound)

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
```

### Settings
The ```ServiceConfig``` structure allow to setup the following parameters:
```golang
type ServiceConfig struct {
    Interface    string           // Host interface to by bind
    Port         int              // Port to listen to
    WriteTimeout time.Duration    // Response Write Timeout
    ReadTimeout  time.Duration    // Request Read Timeout
}
```

### Utility functions

- ```func RespondWithJSONError(w http.ResponseWriter, code int, err error)```: Write to response the parameter error in JSON format
- ```func RespondWithJSONMessage(w http.ResponseWriter, code int, message string)```: Write to response the parameter message in JSON format
- ```func RespondWithJSON(w http.ResponseWriter, code int, payload interface{})```: Write to response the parameter payload as an arbitrary data structure in JSON format
- ```func FixFileName(name string) string```: Return a valid file name representation for the OS file system
- ```func SaveFileFromRequest(r *http.Request, formInputName string, dest string) error```: Save a file sended by the client
- ```func SaveTmpFileFromRequest(r *http.Request, formInputName string, destFolder string) (string, error)```: Save a file sended by the client as a temporal file. Temporal files names include an UID prefix in the format [XXXXXXXX].[REQUEST_FILE_NAME]
- ```func ParseAuthorizationHeader(r *http.Request) string```: Return the value of Authorization hedaer and remove the prefix "Bearer" if present
