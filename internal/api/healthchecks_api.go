package api

import (
	"log"
	"net/http"

	"github.com/SevcikMichal/microfrontends-controller/internal/configuration"
)

// Currently only returns ok to check if the API server is listening
func GetHealthInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("Checking api health")
	w.WriteHeader(http.StatusOK)
}

func CheckHealth(r *http.Request) error {
	_, err := http.Get("http://localhost:" + configuration.GetHttpPort() + "/healthz")
	return err
}
