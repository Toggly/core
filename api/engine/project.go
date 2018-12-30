package engine

import (
	"github.com/Toggly/core/api"
	"github.com/Toggly/core/domain"
	"github.com/Toggly/core/storage"
	"github.com/Toggly/core/util"
)

type projectAPI struct {
	owner  string
	engine *APIEngine
}

func (a *projectAPI) s() storage.ProjectStorage {
	return a.engine.Storage.Projects(a.owner)
}

func (a *projectAPI) List() ([]*domain.Project, error) {
	return a.s().List()
}

func (a *projectAPI) Get(code string) (*domain.Project, error) {
	p, err := a.s().Get(code)
	if err == storage.ErrNotFound {
		return nil, api.ErrProjectNotFound
	}
	return p, err
}

func checkProjectParams(code, description, status string) error {
	if code == "" {
		return api.NewBadRequest("Project code not specified")
	}
	if status != domain.ProjectStatusActive && status != domain.ProjectStatusDisabled {
		return api.NewBadRequest("Project status can be `%s` or `%s`", domain.ProjectStatusActive, domain.ProjectStatusDisabled)
	}
	return nil
}

func (a *projectAPI) Create(info *api.ProjectInfo) (*domain.Project, error) {
	if err := checkProjectParams(info.Code, info.Description, info.Status); err != nil {
		return nil, err
	}
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
	// TODO: create default env
	return newProj, nil
}

func (a *projectAPI) Update(info *api.ProjectInfo) (*domain.Project, error) {
	if err := checkProjectParams(info.Code, info.Description, info.Status); err != nil {
		return nil, err
	}
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

func (a *projectAPI) Delete(code string) error {
	// TODO: delete default environment
	// TODO: check if project not empty
	return a.s().Delete(code)
}
