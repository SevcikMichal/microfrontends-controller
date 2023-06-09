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
