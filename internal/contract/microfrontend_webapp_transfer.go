package contract

type MicroFrontendWebAppTransfer struct {
	MicroFrontendElementTransfer
	Title    string `json:"title"`
	Details  string `json:"details"`
	Path     string `json:"path"`
	Priority int    `json:"priority"`
	Icon     string `json:"icon,omitempty"`
	IsActive bool   `json:"isActive"`
}
