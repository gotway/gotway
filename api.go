package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var (
	address = "0.0.0.0:8080"
	baseurl = "/api"
)

type apiError struct {
	Error   error
	Message string
	Code    int
	Request string
}

func (apiError apiError) LogError() string {
	return apiError.Error.Error() + ";" + apiError.Message + ";" + strconv.Itoa(apiError.Code) + ";" + apiError.Request
}

//APIInit Starts a new HTTP server
func APIInit() {
	router := mux.NewRouter()
	router.HandleFunc(baseurl+"/salute", getSalute).Methods("GET")
	log.Println("Server listening at", address)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getSalute(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(2)%2 == 0 {
		apiError := apiError{errors.New("not found"), "No salute found for this request", 404, r.URL.String()}
		log.Println(apiError.LogError())
		http.Error(w, "No salute for you", http.StatusNotFound)
	} else {
		fmt.Fprintf(w, "Hello\n")
	}
}
