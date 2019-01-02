package domain

import "fmt"

// Parameter types enum
const (
	ParameterTypeBool   = "bool"
	ParameterTypeString = "string"
	ParameterTypeInt    = "int"
)

// Parameter type
type Parameter struct {
	Owner         string        `json:"owner"`
	Code          string        `json:"code"`
	Project       string        `json:"project"`
	Environment   string        `json:"environment"`
	Group         string        `json:"group"`
	Description   string        `json:"description"`
	Type          string        `json:"type"`
	Value         interface{}   `json:"value"`
	AllowedValues []interface{} `json:"allowed_values,omitempty" bson:"allowed_values,omitempty"`
}

// Key returns full group key
func (p *Parameter) Key() string {
	return fmt.Sprintf("owner: %s, project: %s, env: %s, group: %s, parameter: %s", p.Owner, p.Project, p.Environment, p.Group, p.Code)
}
