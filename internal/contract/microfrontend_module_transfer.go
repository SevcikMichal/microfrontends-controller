package contract

type MicroFrontendModuleTransfer struct {
	LoadURL string   `json:"load_url"`
	Styles  []string `json:"styles,omitempty"`
}
