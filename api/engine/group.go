package engine

import (
	"github.com/Toggly/core/api"
	"github.com/Toggly/core/domain"
)

type groupAPI struct {
	owner   string
	project string
	env     string
	engine  *APIEngine
}

func (a *groupAPI) List() ([]*domain.Group, error) {
	return nil, nil
}

func (a *groupAPI) Get(code string) (*domain.Group, error) {
	return nil, nil
}

func (a *groupAPI) Create(info api.GroupInfo) (*domain.Group, error) {
	return nil, nil
}

func (a *groupAPI) Update(info api.GroupInfo) (*domain.Group, error) {
	return nil, nil
}

func (a *groupAPI) Delete(code string) error {
	return nil
}

func (a *groupAPI) Parameters(group string) api.ParameterAPI {
	return &parameterAPI{
		owner:   a.owner,
		project: a.project,
		env:     a.env,
		group:   group,
		engine:  a.engine,
	}
}
