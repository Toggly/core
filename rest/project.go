package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Toggly/core/api"
	"github.com/Toggly/core/domain"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type projectCreateRequest struct {
	Code        string
	Description string
	Status      string
}

type projectRestAPI struct {
	API      api.TogglyAPI
	Log      zerolog.Logger
	LogLevel zerolog.Level
}

func (a *projectRestAPI) Routes() chi.Router {
	router := chi.NewRouter()
	router.Group(func(group chi.Router) {
		group.Get("/", a.list)
		group.Post("/", a.createProject)
		group.Put("/", a.updateProject)
		group.Get("/{project_code}", a.getProject)
		// group.Delete("/{project_code}", a.deleteProject)
	})
	return router
}

func (a *projectRestAPI) engine(r *http.Request) api.ProjectAPI {
	return a.API.Projects(owner(r))
}

func (a *projectRestAPI) list(w http.ResponseWriter, r *http.Request) {
	log := WithRequest(a.Log, r)
	list, err := a.engine(r).List()
	if err != nil {
		log.Error().Err(err).Msg("Can't get projects list")
		ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	JSONResponse(w, r, list)
}

func (a *projectRestAPI) getProject(w http.ResponseWriter, r *http.Request) {
	log := WithRequest(a.Log, r)
	proj, err := a.engine(r).Get(projectCode(r))
	if err != nil {
		log.Error().Err(err).Msg("Can't get project")
	}
	_ = proj
}

func (a *projectRestAPI) createProject(w http.ResponseWriter, r *http.Request) {
	a.createUpdate(w, r, true)
}

func (a *projectRestAPI) updateProject(w http.ResponseWriter, r *http.Request) {
	a.createUpdate(w, r, false)
}

func (a *projectRestAPI) createUpdate(w http.ResponseWriter, r *http.Request, create bool) {
	log := WithRequest(a.Log, r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Can't read request body")
		ErrorResponse(w, r, err, http.StatusInternalServerError)
		return
	}
	proj := &projectCreateRequest{}
	err = json.Unmarshal(body, proj)
	if err != nil {
		log.Error().Err(err).Msg("Can't parse request body")
		ErrorResponse(w, r, errors.New("Bad request"), http.StatusBadRequest)
		return
	}
	info := &api.ProjectInfo{
		Code:        proj.Code,
		Description: proj.Description,
		Status:      proj.Status,
	}
	var p *domain.Project
	if create {
		p, err = a.engine(r).Create(info)
	} else {
		p, err = a.engine(r).Update(info)
	}
	if err != nil {
		log.Error().Err(err).Msg("Can't save/update project")
		return
	}
	JSONResponse(w, r, p)
}
