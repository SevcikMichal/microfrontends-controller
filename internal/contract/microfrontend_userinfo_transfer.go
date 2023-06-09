package contract

type MicroFrontendUserInfoTransfer struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}
