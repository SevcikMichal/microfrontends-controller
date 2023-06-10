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

package model

import (
	"strings"

	microfrontendv1alpha1 "github.com/SevcikMichal/microfrontends-controller/api/v1alpha1"
	"github.com/peteprogrammer/go-automapper"
)

func CreateFrontendConfigFromWebComponent(webComponent *microfrontendv1alpha1.WebComponent) *MicroFrontendConfig {
	frontendConfig := &MicroFrontendConfig{}
	automapper.MapLoose(webComponent.Spec, frontendConfig)
	frontendConfig.MicroFrontendName = webComponent.ObjectMeta.Name
	frontendConfig.MicroFrontendNamespace = webComponent.ObjectMeta.Namespace
	frontendConfig.MicroFrontendLabels = webComponent.ObjectMeta.Labels
	return frontendConfig
}

func RebaseUri(uri string) string {
	base := "/" // TODO: Get from configuration
	baseShort := base
	baseFull := base
	if strings.HasSuffix(base, "/") {
		baseShort = base[:len(base)-1]
	} else {
		baseFull = base + "/"
	}
	if strings.HasPrefix(uri, "/") {
		return baseShort + uri
	} else {
		return baseFull + uri
	}
}
