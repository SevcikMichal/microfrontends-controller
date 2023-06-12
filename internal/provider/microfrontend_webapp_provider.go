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

package provider

import (
	"strings"
	"sync"

	"github.com/SevcikMichal/microfrontends-controller/contract"
	"github.com/SevcikMichal/microfrontends-controller/internal/model"
	"k8s.io/apimachinery/pkg/types"
)

func (r *MicroFrontendProvider) getWebAppTransfers() []*contract.MicroFrontendWebAppTransfer {
	result := []*contract.MicroFrontendWebAppTransfer{}

	webAppTransfersMap, _ := r.MicroFrontendTransferStorage.LoadOrStore(appTransfersKey, &sync.Map{})
	webAppTransferMapCasted := webAppTransfersMap.(*sync.Map)

	// Aggregate all web app transfers
	webAppTransferMapCasted.Range(func(key, value interface{}) bool {
		result = append(result, value.([]*contract.MicroFrontendWebAppTransfer)...)
		return true
	})

	return result
}

func (r *MicroFrontendProvider) updateWebAppTransfers(key types.UID, microFrontendConfig *model.MicroFrontendConfig) {
	// Get a map of all web app transfers that maps Resource UID to all web app transfers that belongs to it
	webAppTransfersMap, _ := r.MicroFrontendTransferStorage.LoadOrStore(appTransfersKey, &sync.Map{})
	webAppTransferMapCasted := webAppTransfersMap.(*sync.Map)

	// Generate a list of all web app transfers that belongs to the given Resource UID
	webAppTransfers := []*contract.MicroFrontendWebAppTransfer{}
	for _, navigation := range microFrontendConfig.Navigations {
		webAppTransfers = append(webAppTransfers, convertFrontendConfigToAppTransfer(microFrontendConfig, navigation))
	}

	// Store the list of all web app transfers that belongs to the given Resource UID
	webAppTransferMapCasted.Store(key, webAppTransfers)

	// Store the map of all web app transfers that maps Resource UID to all web app transfers that belongs to its
	r.MicroFrontendTransferStorage.Store(appTransfersKey, webAppTransferMapCasted)
}

func (r *MicroFrontendProvider) deleteWebAppTransfers(key types.UID) {
	// Get a map of all web app transfers that maps Resource UID to all web app transfers that belongs to it
	webAppTransfersMap, _ := r.MicroFrontendTransferStorage.LoadOrStore(appTransfersKey, &sync.Map{})
	webAppTransferMapCasted := webAppTransfersMap.(*sync.Map)

	// Delete references to this resouce UID
	webAppTransferMapCasted.Delete(key)

	// Store the map of all web app transfers that maps Resource UID to all web app transfers that belongs to its
	r.MicroFrontendTransferStorage.Store(appTransfersKey, webAppTransferMapCasted)
}

func convertFrontendConfigToAppTransfer(frontendConfig *model.MicroFrontendConfig, navigation model.MicroFrontendNavigation) *contract.MicroFrontendWebAppTransfer {
	finalModuleUri := frontendConfig.ExtractModuleUri()

	module := &contract.MicroFrontendModuleTransfer{
		LoadURL: finalModuleUri,
		Styles:  frontendConfig.ExtractStyles(finalModuleUri),
	}

	originalAttributes := navigation.ExtractAttributes()
	extractedAttributes := []*contract.MicroFrontendAttributeTransfer{}
	if len(originalAttributes) > 0 {
		for _, attribute := range originalAttributes {
			extractedAttributes = append(extractedAttributes, attribute.ToContract())
		}
	}

	element := &contract.MicroFrontendElementTransfer{
		MicroFrontendModuleTransfer: module,
		Element:                     navigation.Element,
		Attributes:                  extractedAttributes,
		Labels:                      frontendConfig.ExtractLabels(),
		Roles:                       navigation.ExtractRoles(),
	}

	navigationPath := navigation.Path
	var iconPath string
	if strings.HasSuffix(navigationPath, "/") {
		iconPath = navigationPath
	} else {
		navigationPath = navigationPath + "/"
		iconPath = navigationPath
	}

	if navigation.Icon != nil {
		iconPath = navigation.Icon.ExtractIconPath(iconPath)
	} else {
		iconPath = ""
	}

	webApp := &contract.MicroFrontendWebAppTransfer{
		MicroFrontendElementTransfer: element,
		Title:                        navigation.Title,
		Details:                      navigation.Details,
		Path:                         navigationPath,
		Priority:                     *navigation.Priority,
		Icon:                         iconPath,
	}

	return webApp
}
