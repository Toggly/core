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

// ForOwner api method
func (e *APIEngine) ForOwner(owner string) api.OwnerAPI {
	return &OwnerAPI{
		owner:   owner,
		storage: e.Storage,
		log:     e.Log,
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
