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
	"fmt"
	"net/url"
)

type MicroFrontendConfig struct {
	ModuleUri              string
	Preload                *bool                         `json:"preload,omitempty"`
	Proxy                  *bool                         `json:"proxy,omitempty"`
	HashSuffix             string                        `json:"hash-suffix,omitempty"`
	StyleRelativePaths     []string                      `json:"style-relative-paths,omitempty"`
	ContextElements        []MicroFrontendContextElement `json:"context-elements,omitempty"`
	Navigations            []MicroFrontendNavigation     `json:"navigations,omitempty"`
	MicroFrontendNamespace string                        `json:"microfrontend-namespace,omitempty"`
	MicroFrontendName      string                        `json:"microfrontend-name,omitempty"`
	MicroFrontendLabels    map[string]string             `json:"microfrontend-labels,omitempty"`
}

func (frontendConfig *MicroFrontendConfig) ExtractModuleUri() string {
	if frontendConfig.ModuleUri == "built-in" {
		return ""
	} else if *frontendConfig.Proxy {
		suffix := ""

		if frontendConfig.HashSuffix != "" {
			suffix = "." + frontendConfig.HashSuffix
		}

		webComponentUri := RebaseUri("/web-components/")
		moduleUri := fmt.Sprintf("%s%s/%s/%s%s.jsm", webComponentUri, frontendConfig.MicroFrontendNamespace, frontendConfig.MicroFrontendName, frontendConfig.MicroFrontendName, suffix)
		return moduleUri
	} else {
		return frontendConfig.ModuleUri
	}
}

func (frontendConfig *MicroFrontendConfig) ExtractStyles(finalModuleUri string) []string {
	styles := frontendConfig.StyleRelativePaths

	if len(styles) == 0 {
		return []string{}
	}

	resolvedStyles := make([]string, len(styles))
	baseUrl, _ := url.Parse(finalModuleUri)
	for i, style := range styles {
		relativeUrl, _ := url.Parse(style)
		resolvedStyles[i] = baseUrl.ResolveReference(relativeUrl).String()
	}
	return resolvedStyles
}

func (frontendConfig *MicroFrontendConfig) ExtractLabels() map[string]string {
	labels := map[string]string{}

	if len(frontendConfig.MicroFrontendLabels) > 0 {
		labels = frontendConfig.MicroFrontendLabels
	}

	return labels
}
