package api

import (
	"log"
	"net/http"
)

// Currently only returns ok to check if the API server is listening
func GetHealthInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("Checking api health")
	w.WriteHeader(http.StatusOK)
}

func CheckHealth(r *http.Request) error {
	_, err := http.Get("http://localhost:10000/healthz")
	return err
}
