package engine

import (
	"github.com/Toggly/core/api"
	"github.com/Toggly/core/domain"
)

// ProjectAPI type
type ProjectAPI struct{}

// List projects
func (a *ProjectAPI) List() ([]*domain.Project, error) {
	return nil, nil
}

// Get project info
func (a *ProjectAPI) Get(code string) (*domain.Project, error) {
	return nil, nil
}

// Create project
func (a *ProjectAPI) Create(info *api.ProjectInfo) (*domain.Project, error) {
	return nil, nil
}

// Update project
func (a *ProjectAPI) Update(info *api.ProjectInfo) (*domain.Project, error) {
	return nil, nil
}

// Delete project
func (a *ProjectAPI) Delete(code string) error {
	return nil
}
