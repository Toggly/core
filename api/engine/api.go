package engine

import (
	"github.com/Toggly/core/api"
	"github.com/Toggly/core/storage"
	"github.com/rs/zerolog"
)

// NewTogglyAPI returns api engine
func NewTogglyAPI(storage storage.DataStorage, log zerolog.Logger) api.TogglyAPI {
	return &Engine{
		storage: storage,
		log:     log,
	}
}

// Engine type
type Engine struct {
	storage storage.DataStorage
	log     zerolog.Logger
}

// ForOwner returns owner api
func (e *Engine) ForOwner(owner string) api.OwnerAPI {
	return &OwnerAPI{
		owner:   owner,
		storage: e.storage,
		log:     e.log,
	}
}

// OwnerAPI type
type OwnerAPI struct {
	owner   string
	storage storage.DataStorage
	log     zerolog.Logger
}

// Projects returns project api
func (o *OwnerAPI) Projects() api.ProjectAPI {
	return &ProjectAPI{*o}
}
