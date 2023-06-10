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

// TODO: Improve the implementation of MicroFrontendProvider currently it is not very sofiticated

import (
	"sync"

	"github.com/SevcikMichal/microfrontends-controller/contract"
	"github.com/SevcikMichal/microfrontends-controller/internal/model"
	"k8s.io/apimachinery/pkg/types"
)

const (
	appTransfersKey     = "app-transfers"
	preloadTrnasfersKey = "preload-transfers"
	contextTransfersKey = "context-transfers"
)

type MicroFrontendProvider struct {
	MicroFrontendTransferStorage *sync.Map
}

func (r *MicroFrontendProvider) SetMicroFrontendConfig(key types.UID, microFrontendConfig *model.MicroFrontendConfig) {
	r.updateWebAppTransfers(key, microFrontendConfig)
	r.updatePreloadTransfers(key, microFrontendConfig)
	r.updateContextTransfers(key, microFrontendConfig)
}

func (r *MicroFrontendProvider) DeleteMicroFrontendConfig(key types.UID) {
	r.deleteWebAppTransfers(key)
	r.deletePreloadTransfers(key)
	r.deleteContextTransfers(key)
}

func (r *MicroFrontendProvider) GetMicroFrontendConfigTransfer() *contract.MicroFrontendConfigurationTransfer {
	frontendConfigTransfer := &contract.MicroFrontendConfigurationTransfer{
		Apps:     r.getWebAppTransfers(),
		Preload:  r.getPreloadTransfers(),
		Contexts: r.getContextTransfers(),
		// Anonymous: ,
		// User: ,
	}

	return frontendConfigTransfer
}

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

func (r *MicroFrontendProvider) getPreloadTransfers() []*contract.MicroFrontendModuleTransfer {
	result := []*contract.MicroFrontendModuleTransfer{}

	preloadTransfersMap, _ := r.MicroFrontendTransferStorage.LoadOrStore(preloadTrnasfersKey, &sync.Map{})
	preloadTransferMapCasted := preloadTransfersMap.(*sync.Map)

	// Aggregate all preload transfers
	preloadTransferMapCasted.Range(func(key, value interface{}) bool {
		result = append(result, value.(*contract.MicroFrontendModuleTransfer))
		return true
	})

	return result
}

func (r *MicroFrontendProvider) getContextTransfers() []*contract.MicroFrontendContextTransfer {
	result := []*contract.MicroFrontendContextTransfer{}

	contextTransfersMap, _ := r.MicroFrontendTransferStorage.LoadOrStore(contextTransfersKey, &sync.Map{})
	contextTransferMapCasted := contextTransfersMap.(*sync.Map)

	// Aggregate all context transfers
	contextTransferMapCasted.Range(func(key, value interface{}) bool {
		result = append(result, value.([]*contract.MicroFrontendContextTransfer)...)
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

func (r *MicroFrontendProvider) updatePreloadTransfers(key types.UID, microFrontendConfig *model.MicroFrontendConfig) {
	// if preload is false we don't want to add the preload for this resource (if it is and update we want to remove it)
	if !*microFrontendConfig.Preload {
		r.deletePreloadTransfers(key)
		return
	}

	// Get a map of all preload transfers that maps Resource UID to all preload transfers that belongs to it
	preloadTransfersMap, _ := r.MicroFrontendTransferStorage.LoadOrStore(preloadTrnasfersKey, &sync.Map{})
	preloadTransferMapCasted := preloadTransfersMap.(*sync.Map)

	// Create preload transfer from the resource
	preloadTransfer := convertFrontendConfigToPreloadTransfer(microFrontendConfig)

	// Store the list of all preload transfers that belongs to the given Resource UID
	preloadTransferMapCasted.Store(key, preloadTransfer)

	// Store the map of all preload transfers that maps Resource UID to all preload transfers that belongs to its
	r.MicroFrontendTransferStorage.Store(preloadTrnasfersKey, preloadTransferMapCasted)
}

func (r *MicroFrontendProvider) deletePreloadTransfers(key types.UID) {
	// Get a map of all preload transfers that maps Resource UID to all preload transfers that belongs to it
	preloadTransfersMap, _ := r.MicroFrontendTransferStorage.LoadOrStore(preloadTrnasfersKey, &sync.Map{})
	preloadTransferMapCasted := preloadTransfersMap.(*sync.Map)

	// Delete references to this resouce UID
	preloadTransferMapCasted.Delete(key)

	// Store the map of all preload transfers that maps Resource UID to all preload transfers that belongs to its
	r.MicroFrontendTransferStorage.Store(preloadTrnasfersKey, preloadTransferMapCasted)
}

func (r *MicroFrontendProvider) updateContextTransfers(key types.UID, microFrontendConfig *model.MicroFrontendConfig) {
	// Get a map of all context transfers that maps Resource UID to all context transfers that belongs to it
	contextTransfersMap, _ := r.MicroFrontendTransferStorage.LoadOrStore(contextTransfersKey, &sync.Map{})
	contextTransferMapCasted := contextTransfersMap.(*sync.Map)

	// Generate a list of all context transfers that belongs to the given Resource UID
	contextTransfers := []*contract.MicroFrontendContextTransfer{}
	for _, context := range microFrontendConfig.ContextElements {
		contextTransfers = append(contextTransfers, convertFrontendConfigToContextTransfer(microFrontendConfig, context))
	}

	// Store the list of all context transfers that belongs to the given Resource UID
	contextTransferMapCasted.Store(key, contextTransfers)

	// Store the map of all context transfers that maps Resource UID to all context transfers that belongs to its
	r.MicroFrontendTransferStorage.Store(contextTransfersKey, contextTransferMapCasted)
}

func (r *MicroFrontendProvider) deleteContextTransfers(key types.UID) {
	// Get a map of all context transfers that maps Resource UID to all context transfers that belongs to it
	contextTransfersMap, _ := r.MicroFrontendTransferStorage.LoadOrStore(contextTransfersKey, &sync.Map{})
	contextTransferMapCasted := contextTransfersMap.(*sync.Map)

	// Delete references to this resouce UID
	contextTransferMapCasted.Delete(key)

	// Store the map of all context transfers that maps Resource UID to all context transfers that belongs to its
	r.MicroFrontendTransferStorage.Store(contextTransfersKey, contextTransferMapCasted)
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

	webApp := &contract.MicroFrontendWebAppTransfer{
		MicroFrontendElementTransfer: element,
		Title:                        navigation.Title,
		Details:                      navigation.Details,
		Path:                         navigation.Path,
		Priority:                     *navigation.Priority,
		// Icon:                         navigation.Icon.Url, TODO: What is happening here in prolog?
	}

	return webApp
}

func convertFrontendConfigToPreloadTransfer(frontendConfig *model.MicroFrontendConfig) *contract.MicroFrontendModuleTransfer {
	finalModuleUri := frontendConfig.ExtractModuleUri()

	module := &contract.MicroFrontendModuleTransfer{
		LoadURL: finalModuleUri,
		Styles:  frontendConfig.ExtractStyles(finalModuleUri),
	}

	return module
}

func convertFrontendConfigToContextTransfer(frontendConfig *model.MicroFrontendConfig, context model.MicroFrontendContextElement) *contract.MicroFrontendContextTransfer {
	finalModuleUri := frontendConfig.ExtractModuleUri()

	module := &contract.MicroFrontendModuleTransfer{
		LoadURL: finalModuleUri,
		Styles:  frontendConfig.ExtractStyles(finalModuleUri),
	}

	originalAttributes := context.ExtractAttributes()
	extractedAttributes := []*contract.MicroFrontendAttributeTransfer{}
	if len(originalAttributes) > 0 {
		for _, attribute := range originalAttributes {
			extractedAttributes = append(extractedAttributes, attribute.ToContract())
		}
	}

	element := &contract.MicroFrontendElementTransfer{
		MicroFrontendModuleTransfer: module,
		Element:                     context.Element,
		Attributes:                  extractedAttributes,
		Labels:                      frontendConfig.ExtractLabels(),
		Roles:                       context.ExtractRoles(),
	}

	contextTransfer := &contract.MicroFrontendContextTransfer{
		MicroFrontendElementTransfer: element,
		ContextNames:                 context.ExtractContextNames(),
	}

	return contextTransfer
}
