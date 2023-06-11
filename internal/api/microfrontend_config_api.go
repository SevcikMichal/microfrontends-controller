/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/SevcikMichal/microfrontends-controller/contract"
	"github.com/SevcikMichal/microfrontends-controller/internal/configuration"
	"github.com/SevcikMichal/microfrontends-controller/internal/provider"
)

type MicroFrontendConfigApi struct {
	MicroFrontendProvider *provider.MicroFrontendProvider
}

func (api *MicroFrontendConfigApi) GetMicroFrontendConfigs(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get micro frontend configs started.")

	frontendConfig := api.MicroFrontendProvider.GetMicroFrontendConfigTransfer()

	userContext := getUserContext(&r.Header)
	if userContext != nil {
		frontendConfig.User = userContext
	} else {
		frontendConfig.Anonymous = new(bool)
		*frontendConfig.Anonymous = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(frontendConfig)
}

func getUserContext(r *http.Header) *contract.MicroFrontendUserInfoTransfer {
	userId := r.Get(configuration.GetUserIdHeader())
	userEmail := r.Get(configuration.GetUserEmailHeader())
	userName := r.Get(configuration.GetUserNameHeader())
	userRoles := r.Get(configuration.GetUserRolesHeader())

	if len(userId) == 0 || len(userEmail) == 0 || len(userName) == 0 || len(userRoles) == 0 {
		return nil
	}

	userRolesArray := strings.Split(userRoles, ",")

	return &contract.MicroFrontendUserInfoTransfer{
		ID:    userId,
		Email: userEmail,
		Name:  userName,
		Roles: userRolesArray,
	}
}
