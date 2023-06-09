package api

import (
	"encoding/json"
	"log"
	"net/http"

	"sync"

	"github.com/SevcikMichal/microfrontends-controller/contract"
	"github.com/SevcikMichal/microfrontends-controller/internal/model"
)

type MicroFrontendConfigApi struct {
	MicroFrontendConfigs *sync.Map
}

func (api *MicroFrontendConfigApi) GetMicroFrontendConfigs(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to get micro frontend configs started.")

	microFrontendWebAppTransfers := []*contract.MicroFrontendWebAppTransfer{}
	//microFrontendContextTransfer = []*contract.MicroFrontendContextTransfer{}
	//microFrontendModuleTransfer = []*contract.MicroFrontendModuleTransfer{}

	api.MicroFrontendConfigs.Range(func(key, microFrontendConfig interface{}) bool {
		config, ok := microFrontendConfig.(*model.MicroFrontendConfig)

		if !ok {
			log.Fatal("Could not convert micro frontend config to model.MicroFrontendConfig.")
		}

		for _, navigation := range config.Navigations {
			module := contract.MicroFrontendModuleTransfer{
				LoadURL: config.ModuleUri,
				//Styles:  config.StyleRelativePaths,
			}

			// Create a new MicroFrontendElementAttributeTransfer object
			// attributes := []contract.MicroFrontendElementAttributeTransfer{}
			// for _, attribute := range navigation.Attributes {
			// 	attributes = append(attributes, contract.MicroFrontendElementAttributeTransfer{
			// 		Name:  attribute.Name,
			// 		Value: attribute.Value,
			// 	})
			// }

			// Create a new MicroFrontendElementTransfer object
			element := contract.MicroFrontendElementTransfer{
				MicroFrontendModuleTransfer: module,
				Element:                     navigation.Element,
				//Attributes:                  attributes,
				//Labels:                      map[string]string{"label1": "value1", "label2": "value2"},
				//Roles:                       navigation.Roles,
			}

			// Create a new MicroFrontendWebAppTransfer object
			webApp := &contract.MicroFrontendWebAppTransfer{
				MicroFrontendElementTransfer: element,
				Title:                        navigation.Title,
				Details:                      navigation.Details,
				Path:                         navigation.Path,
				//Priority:                     navigation.Priority,
				//Icon:                         navigation.Icon.Url,
			}

			microFrontendWebAppTransfers = append(microFrontendWebAppTransfers, webApp)
		}

		return true
	})

	json.NewEncoder(w).Encode(microFrontendWebAppTransfers)
}
