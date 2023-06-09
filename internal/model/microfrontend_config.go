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

type MicroFrontendConfig struct {
	ModuleUri          string
	Preload            bool                          `json:"preload,omitempty"`
	Proxy              bool                          `json:"proxy,omitempty"`
	HashSuffix         string                        `json:"hash-suffix,omitempty"`
	StyleRelativePaths []string                      `json:"style-relative-paths,omitempty"`
	ContextElements    []MicroFrontendContextElement `json:"context-elements,omitempty"`
	Navigations        []MicroFrontendNavigation     `json:"navigations,omitempty"`
}
