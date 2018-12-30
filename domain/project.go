package domain

import (
	"fmt"
	"time"
)

// ProjectStatus enum
const (
	ProjectStatusActive   = "active"
	ProjectStatusDisabled = "disabled"
)

// Project type
type Project struct {
	Code        string    `json:"code"`
	Owner       string    `json:"owner"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	RegDate     time.Time `json:"reg_date" bson:"reg_date"`
}

// Key returns project full key
func (p *Project) Key() string {
	return fmt.Sprintf("owner: %s, project: %s", p.Owner, p.Code)
}
