package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SevcikMichal/microfrontends-controller/internal/provider"
)

type MicroFrontendConfigApi struct {
	MicroFrontendProvider *provider.MicroFrontendProvider
}

func (api *MicroFrontendConfigApi) GetMicroFrontendConfigs(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get micro frontend configs started.")

	json.NewEncoder(w).Encode(api.MicroFrontendProvider.GetMicroFrontendConfigTransfer())
}
