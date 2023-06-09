package model

type MicroFrontendContextElement struct {
	ContextNames []string                 `json:"context-names"`
	Element      string                   `json:"element"`
	Priority     int                      `json:"priority,omitempty"`
	Attributes   []MicroFrontendAttribute `json:"attributes,omitempty"`
	Roles        []string                 `json:"roles,omitempty"`
}
