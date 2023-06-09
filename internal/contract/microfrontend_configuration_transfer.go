package contract

type MicroFrontendConfigurationTransfer struct {
	Preload   []MicroFrontendModuleTransfer  `json:"preload"`
	Apps      []MicroFrontendWebAppTransfer  `json:"apps"`
	Contexts  []MicroFrontendContextTransfer `json:"contexts"`
	Anonymous bool                           `json:"anonymous,omitempty"`
	User      *MicroFrontendUserInfoTransfer `json:"user,omitempty"`
}
