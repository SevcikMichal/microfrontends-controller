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

type MicroFrontendContextElement struct {
	ContextNames []string                 `json:"context-names"`
	Element      string                   `json:"element"`
	Priority     *int                     `json:"priority,omitempty"`
	Attributes   []MicroFrontendAttribute `json:"attributes,omitempty"`
	Roles        []string                 `json:"roles,omitempty"`
}

// TODO: Check wheteher some string spliting should be still applied here or not
func (context *MicroFrontendContextElement) ExtractRoles() []string {
	if len(context.Roles) > 0 {
		return context.Roles
	} else {
		return []string{"*"}
	}
}

func (context *MicroFrontendContextElement) ExtractAttributes() []MicroFrontendAttribute {
	attributes := []MicroFrontendAttribute{}

	if len(context.Attributes) > 0 {
		attributes = context.Attributes
	}

	return attributes
}

func (context *MicroFrontendContextElement) ExtractContextNames() []string {
	contextNames := []string{}

	if len(context.ContextNames) > 0 {
		contextNames = context.ContextNames
	}

	return contextNames
}
