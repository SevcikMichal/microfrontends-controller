package contract

type MicroFrontendElementTransfer struct {
	MicroFrontendModuleTransfer
	Element    string                                  `json:"element"`
	Attributes []MicroFrontendElementAttributeTransfer `json:"attributes"`
	Labels     map[string]string                       `json:"labels,omitempty"`
	Roles      []string                                `json:"roles,omitempty"`
}
