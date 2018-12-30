package engine

import (
	"github.com/Toggly/core/api"
	"github.com/Toggly/core/storage"
	"github.com/rs/zerolog"
)

// APIEngine type
type APIEngine struct {
	Storage storage.DataStorage
	Log     zerolog.Logger
}

// Projects returns project api
func (a *APIEngine) Projects(owner string) api.ProjectAPI {
	return &projectAPI{
		owner:  owner,
		engine: a,
	}
}

// Environments returns environments api
func (a *APIEngine) Environments(owner, project string) api.EnvironmentAPI {
	return &environmentAPI{
		owner:   owner,
		project: project,
		engine:  a,
	}
}

// Groups returns groups api
func (a *APIEngine) Groups(owner, project, env string) api.GroupAPI {
	return &groupAPI{
		owner:   owner,
		project: project,
		env:     env,
		engine:  a,
	}
}
