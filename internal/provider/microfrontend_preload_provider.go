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

func convertFrontendConfigToPreloadTransfer(frontendConfig *model.MicroFrontendConfig) *contract.MicroFrontendModuleTransfer {
	finalModuleUri := frontendConfig.ExtractModuleUri()

	module := &contract.MicroFrontendModuleTransfer{
		LoadURL: finalModuleUri,
		Styles:  frontendConfig.ExtractStyles(finalModuleUri),
	}

	return module
}
