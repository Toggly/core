package engine

import (
	"github.com/Toggly/core/api"
	"github.com/Toggly/core/domain"
)

type parameterAPI struct {
	owner   string
	project string
	env     string
	group   string
	engine  *APIEngine
}

func (a *parameterAPI) List() ([]*domain.Parameter, error) {
	return nil, nil
}

func (a *parameterAPI) Get(code string) (*domain.Parameter, error) {
	return nil, nil
}

func (a *parameterAPI) GetBatch(code ...string) ([]*domain.Parameter, error) {
	return nil, nil
}

func (a *parameterAPI) Create(param *api.ParameterInfo) (*domain.Parameter, error) {
	return nil, nil
}

func (a *parameterAPI) Update(param *api.ParameterInfo) (*domain.Parameter, error) {
	return nil, nil
}

func (a *parameterAPI) Delete(code string) error {
	return nil
}
