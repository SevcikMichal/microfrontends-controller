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
