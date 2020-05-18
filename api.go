package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type apiError struct {
	Error   error
	Message string
	Code    int
	Request string
}

var logger *zap.SugaredLogger

func (apiError apiError) LogError() string {
	return apiError.Error.Error() + ";" + apiError.Message + ";" + strconv.Itoa(apiError.Code) + ";" + apiError.Request
}

//NewAPI Starts a new HTTP server
func NewAPI() {
	initializeLogger()
	router := mux.NewRouter()
	router.HandleFunc("/api/salute", saluteHandler).Methods("GET")
	router.HandleFunc("/api/health", healthHandler).Methods("GET")
	port := getEnv("PORT", ":8080")
	logger.Info("Server listening on port ", port)
	logger.Fatal(http.ListenAndServe(port, router))
}

func initializeLogger() {
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()
	logger = zapLogger.Sugar()
}

func saluteHandler(w http.ResponseWriter, r *http.Request) {
	if rand.Intn(2)%2 == 0 {
		apiError := apiError{errors.New("not found"), "No salute found for this request", 404, r.URL.String()}
		logger.Error(apiError.LogError())
		http.Error(w, "No salute for you", http.StatusNotFound)
	} else {
		fmt.Fprintf(w, "Hello\n")
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
