package engine

import (
	"fmt"

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

func (a *ProjectAPI) List() ([]*domain.Project, error) {
	return a.s().List()
}

func (a *ProjectAPI) Get(code string) (*domain.Project, error) {
	p, err := a.s().Get(code)
	if err == storage.ErrNotFound {
		return nil, api.ErrProjectNotFound
	}
	return p, err
}

func checkProjectParams(code, description, status string) error {
	if code == "" {
		return &api.ErrBadRequest{
			Description: "Project code not specified",
		}
	}
	if status != domain.ProjectStatusActive && status != domain.ProjectStatusDisabled {
		return &api.ErrBadRequest{
			Description: fmt.Sprintf("Project status can be `%s` or `%s`", domain.ProjectStatusActive, domain.ProjectStatusDisabled),
		}
	}
	return nil
}

func (a *ProjectAPI) Create(info *api.ProjectInfo) (*domain.Project, error) {
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

func (a *ProjectAPI) Update(info *api.ProjectInfo) (*domain.Project, error) {
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

func (a *ProjectAPI) Delete(code string) error {
	return a.s().Delete(code)
}

func (a *ProjectAPI) For(code string) api.ForProjectAPI {
	return &forProjectAPI{}
}

type forProjectAPI struct{}

func (a *forProjectAPI) Environments() api.EnvironmentAPI {
	return &environmentAPI{}
}
