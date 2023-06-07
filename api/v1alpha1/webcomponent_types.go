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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WebComponentSpec defines the desired state of WebComponent
type WebComponentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// The URI from which the module shall be accessed. The actual module is cached by the controller to improve performance and avoid CORS issues.
	// +kubebuilder:validation:Format=url
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ModuleUri string `json:"module-uri"`

	// The modules are not preloaded by default but only when navigating to some of the subpaths mentioned in the 'navigation' list. Setting this property to true ensures that the module is loaded when the application starts.
	// +kubebuilder:default=true
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Preload bool `json:"preload,omitempty"`

	// This specifies whether the loading of web components should be proxied by the controller. This is useful if the web component is served from within the cluster and cannot be accessed from outside the cluster network. The module will be served from the URL base_controller_url/web-components/web_component_name.jsm. This is the recommended approach for the standard assumed use-case.
	// +kubebuilder:default=false
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Proxy bool `json:"proxy,omitempty"`

	// A hash string used to identify a specific version of the module URI when the controller is proxying it. If the proxy property is set and the hash property is set as well, the final module file name will be web_component_name.hash_suffix.jsm, and the resource will be assumed to never expire. To refresh user agents' caches, the hash value needs to be changed to a new unique value.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	HashSuffix string `json:"hash-suffix,omitempty"`

	// An array of relative paths (relative to module-uri) that contains the CSS styles required for this web component module. Ideally, the styles are either embedded in or loaded by the module. However, certain legacy styles may require an additional link element.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	StyleRelativePaths []string `json:"style-relative-paths,omitempty"`

	// These are components that can be displayed in a special context, such as ufe-app-shell for the top-level application shell or my-menu-item for components to be displayed in a custom menu, and so on.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ContextElements []ContextElement `json:"context-elements,omitempty"`

	// Components with the navigation specification may be used as sub-paths and are considered as workspaces or applications on their own within the composed application shell.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Navigations []Navigation `json:"navigations,omitempty"`
}

type ContextElement struct {

	// This is a list of context names in which this element is intended to be shown.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	ContextNames []string `json:"context-names"`

	// This is the HTML element tag name to use when navigating to the specific path.
	// +kubebuilder:example="my-menu-item"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Element string `json:"element"`

	// This indicates the priority of the navigation entry in lists. Entries with higher priority are displayed before entries with lower priorities, if there is an ordering supported by the list. The default priority is 0.
	// +kubebuilder:default=0
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Priority int `json:"priority,omitempty"`

	// This is a list of key-value pairs that allows you to assign specific attributes to the element. The name field is used as the attribute name, while the value field can be any valid JSON type.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Attributes []Attribute `json:"attributes,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Roles []string `json:"roles,omitempty"`
}

type Attribute struct {

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Name string `json:"name"`

	// +kubebuilder:validation:XPreserveUnknownFields
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Value string `json:"value"`
}

type Navigation struct {

	// By navigating to the specific subpath, the app shell will place the element on the main workspace (content) of the shell.
	// +kubebuilder:example="/my-menu-item"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Path string `json:"path"`

	// The title is used to present links to the particular workspace or to display it in navigation lists, or as a title when on the specific path.
	// +kubebuilder:example="My Menu Item"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Title string `json:"title"`

	// This indicates the priority of the navigation entry in lists. Entries with higher priority are displayed before entries with lower priorities, if there is an ordering supported by the list. The default priority is 0.
	// +kubebuilder:default=0
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Priority int `json:"priority,omitempty"`

	// Additional description is provided to explain the purpose of the component to the user. It is displayed in addition to the title in the navigation lists.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Details string `json:"details,omitempty"`

	// The HTML element tag name to be used when navigating to the specific path.
	// +kubebuilder:example="my-menu-item"
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Element string `json:"element"`

	// This is a list of key-value pairs that allows you to assign specific attributes to the element. The name field is used as the attribute name, while the value field can be any valid JSON type.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Attributes []Attribute `json:"attributes,omitempty"`

	// The optional icon is associated with the navigable component. SVG format is preferred. Either the data property should provide base64 encoded icon/image data or the url to the image source should be specified. The mime property must specify the proper MIME type of the icon/image.
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Icon Icon `json:"icon,omitempty"`

	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Roles []string `json:"roles,omitempty"`
}

type Icon struct {
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Mime string `json:"mime"`

	// +kubebuilder:Validation:Format=byte
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Data string `json:"data"`

	// +kubebuilder:Validation:Format=url
	// +operator-sdk:csv:customresourcedefinitions:type=spec
	Url string `json:"url"`
}

// WebComponentStatus defines the observed state of WebComponent
type WebComponentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// WebComponent is the Schema for the webcomponents API
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=webc
type WebComponent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebComponentSpec   `json:"spec"`
	Status WebComponentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WebComponentList contains a list of WebComponent
type WebComponentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebComponent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebComponent{}, &WebComponentList{})
}
