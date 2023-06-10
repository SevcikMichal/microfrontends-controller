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

type MicroFrontendNavigation struct {
	Path       string                   `json:"path"`
	Title      string                   `json:"title"`
	Priority   int                      `json:"priority,omitempty"`
	Details    string                   `json:"details,omitempty"`
	Element    string                   `json:"element"`
	Attributes []MicroFrontendAttribute `json:"attributes,omitempty"`
	Icon       MicroFrontendIcon        `json:"icon,omitempty"`
	Roles      []string                 `json:"roles,omitempty"`
}

// TODO: Check wheteher some string spliting should be still applied here or not
func (navigation *MicroFrontendNavigation) ExtractRoles() []string {
	roles := []string{}

	if len(navigation.Roles) > 0 {
		roles = navigation.Roles
	}

	return roles
}

func (navigation *MicroFrontendNavigation) ExtractAttributes() []MicroFrontendAttribute {
	attributes := []MicroFrontendAttribute{}

	if len(navigation.Attributes) > 0 {
		attributes = navigation.Attributes
	}

	return attributes
}
