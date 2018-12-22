package engine

import (
	"github.com/Toggly/core/api"
	"github.com/Toggly/core/storage"
	"github.com/rs/zerolog"
)

// NewTogglyAPI returns api engine
func NewTogglyAPI(storage storage.DataStorage, log zerolog.Logger) api.TogglyAPI {
	return &engine{
		storage: storage,
		log:     log,
	}
}

type engine struct {
	storage storage.DataStorage
	log     zerolog.Logger
}

func (e *engine) ForOwner(owner string) api.OwnerAPI {
	return &ownerAPI{
		owner:   owner,
		storage: e.storage,
		log:     e.log,
	}
}

type ownerAPI struct {
	owner   string
	storage storage.DataStorage
	log     zerolog.Logger
}

func (o *ownerAPI) Projects() api.ProjectAPI {
	return &projectAPI{*o}
}
