package engine

import (
	"github.com/Toggly/core/api"
	"github.com/Toggly/core/domain"
	"github.com/Toggly/core/storage"
	"github.com/Toggly/core/util"
)

// ProjectAPI type
type ProjectAPI struct {
	OwnerAPI
}

func (a *ProjectAPI) s() storage.ProjectStorage {
	return a.storage.ForOwner(a.owner).Projects()
}

// List projects
func (a *ProjectAPI) List() ([]*domain.Project, error) {
	return a.s().List()
}

// Get project info
func (a *ProjectAPI) Get(code string) (*domain.Project, error) {
	return a.s().Get(code)
}

// Create project
func (a *ProjectAPI) Create(info *api.ProjectInfo) (*domain.Project, error) {
	newProj := &domain.Project{
		Code:        info.Code,
		Description: info.Description,
		Owner:       a.owner,
		RegDate:     util.Now(),
		Status:      info.Status,
	}
	if err := a.s().Save(newProj); err != nil {
		return nil, err
	}
	return newProj, nil
}

// Update project
func (a *ProjectAPI) Update(info *api.ProjectInfo) (*domain.Project, error) {
	proj, err := a.s().Get(info.Code)
	if err != nil {
		return nil, err
	}
	newProj := &domain.Project{
		Code:        info.Code,
		Description: info.Description,
		Owner:       a.owner,
		RegDate:     proj.RegDate,
		Status:      info.Status,
	}
	err = a.s().Update(newProj)
	if err != nil {
		return nil, err
	}
	return newProj, nil
}

// Delete project
func (a *ProjectAPI) Delete(code string) error {
	return a.s().Delete(code)
}
