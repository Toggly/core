package engine

import (
	"github.com/Toggly/core/api"
	"github.com/Toggly/core/domain"
)

type environmentAPI struct {
}

func (a *environmentAPI) List() ([]*domain.Environment, error) {
	return nil, nil
}

func (a *environmentAPI) Get(code string) (*domain.Environment, error) {
	return nil, nil
}

func (a *environmentAPI) Create(info *api.EnvironmentInfo) (*domain.Environment, error) {
	return nil, nil
}

func (a *environmentAPI) Update(info *api.EnvironmentInfo) (*domain.Environment, error) {
	return nil, nil
}

func (a *environmentAPI) Delete(code string) error {
	return nil
}

func (a *environmentAPI) For(code string) api.ForEnvironmentAPI {
	return nil
}
