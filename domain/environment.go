package domain

import (
	"fmt"
	"time"
)

// Environment type
type Environment struct {
	Code        string    `json:"code"`
	Owner       string    `json:"owner"`
	Project     string    `json:"project"`
	Description string    `json:"description"`
	Protected   bool      `json:"protected"`
	RegDate     time.Time `json:"reg_date" bson:"reg_date"`
}

// Key returns full environment key
func (e *Environment) Key() string {
	return fmt.Sprintf("owner: %s, project: %s, env: %s", e.Owner, e.Project, e.Code)
}
