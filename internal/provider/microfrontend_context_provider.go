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
	"sync"

	"github.com/SevcikMichal/microfrontends-controller/contract"
	"github.com/SevcikMichal/microfrontends-controller/internal/model"
	"k8s.io/apimachinery/pkg/types"
)

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
