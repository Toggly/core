package domain

import "fmt"

// Group type
type Group struct {
	Owner       string `json:"owner"`
	Project     string `json:"project"`
	Environment string `json:"environment"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// Key returns full group key
func (g *Group) Key() string {
	return fmt.Sprintf("owner: %s, project: %s, env: %s, group: %s", g.Owner, g.Project, g.Environment, g.Code)
}
