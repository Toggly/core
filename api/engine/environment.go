package engine

import (
	"github.com/Toggly/core/domain"
)

type environmentAPI struct{}

func (a *environmentAPI) List() ([]*domain.Environment, error) {
	return nil, nil
}

func (a *environmentAPI) Get(code string) (*domain.Environment, error) {
	return nil, nil
}
