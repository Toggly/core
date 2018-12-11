package api

import "github.com/Toggly/core/domain"

// TogglyAPI interface
type TogglyAPI interface{}

// ProjectInfo type
type ProjectInfo struct {
	Code        string
	Description string
	Status      string
}

// ProjectAPI interface
type ProjectAPI interface {
	List() ([]*domain.Project, error)
	Get(code string) (*domain.Project, error)
	Create(info *ProjectInfo) (*domain.Project, error)
	Update(info *ProjectInfo) (*domain.Project, error)
	Delete(code string) error
	Some()
}
