package domain

// Group type
type Group struct {
	Code        string `json:"code"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Project     string `json:"project"`
	Environment string `json:"environment"`
}
