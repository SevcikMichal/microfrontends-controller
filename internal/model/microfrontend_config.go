package model

import "k8s.io/apimachinery/pkg/runtime"

type MicroFrontendConfig struct {
	ModuleUri          string
	Preload            bool                          `json:"preload,omitempty"`
	Proxy              bool                          `json:"proxy,omitempty"`
	HashSuffix         string                        `json:"hash-suffix,omitempty"`
	StyleRelativePaths []string                      `json:"style-relative-paths,omitempty"`
	ContextElements    []MicroFrontendContextElement `json:"context-elements,omitempty"`
	Navigations        []MicroFrontendNavigation     `json:"navigations,omitempty"`
}

type MicroFrontendContextElement struct {
	ContextNames []string                 `json:"context-names"`
	Element      string                   `json:"element"`
	Priority     int                      `json:"priority,omitempty"`
	Attributes   []MicroFrontendAttribute `json:"attributes,omitempty"`
	Roles        []string                 `json:"roles,omitempty"`
}

type MicroFrontendAttribute struct {
	Name  string               `json:"name"`
	Value runtime.RawExtension `json:"value"`
}

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

type MicroFrontendIcon struct {
	Mime string `json:"mime"`
	Data string `json:"data"`
	Url  string `json:"url"`
}
