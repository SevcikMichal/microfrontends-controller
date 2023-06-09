package contract

type MicroFrontendContextTransfer struct {
	MicroFrontendElementTransfer
	ContextNames []string `json:"contextNames"`
}
