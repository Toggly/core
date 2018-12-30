package api

import (
	"errors"
	"fmt"

	"github.com/Toggly/core/domain"
)

const (
	// DefaultEnvName created by default with project. Can't be deleted. Automatically deleted with project.
	DefaultEnvName = "base"
)

var (
	// ErrProjectNotFound error
	ErrProjectNotFound = errors.New("Project not found")
	// ErrProjectNotEmpty error
	ErrProjectNotEmpty = errors.New("Project not empty")
)

// ErrBadRequest type
type ErrBadRequest struct {
	Description string
}

func (e *ErrBadRequest) Error() string {
	return fmt.Sprintf("Bad request: %s", e.Description)
}

// NewBadRequest returns BadRequest object
func NewBadRequest(format string, a ...interface{}) *ErrBadRequest {
	return &ErrBadRequest{Description: fmt.Sprintf(format, a...)}
}

// TogglyAPI interface
type TogglyAPI interface {
	Projects(owner string) ProjectAPI
	Environments(owner, project string) EnvironmentAPI
	Groups(owner, project, env string) GroupAPI
}

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
}

// EnvironmentInfo type
type EnvironmentInfo struct{}

// EnvironmentAPI interface
type EnvironmentAPI interface {
	List() ([]*domain.Environment, error)
	Get(code string) (*domain.Environment, error)
	Create(info *EnvironmentInfo) (*domain.Environment, error)
	Update(info *EnvironmentInfo) (*domain.Environment, error)
	Delete(code string) error
	Parameters(env string) ParameterAPI
}

// GroupInfo type
type GroupInfo struct{}

// GroupAPI interface
type GroupAPI interface {
	List() ([]*domain.Group, error)
	Get(code string) (*domain.Group, error)
	Create(info GroupInfo) (*domain.Group, error)
	Update(info GroupInfo) (*domain.Group, error)
	Delete(code string) error
	Parameters(group string) ParameterAPI
}

// ParameterInfo type
type ParameterInfo struct{}

// ParameterAPI interface
type ParameterAPI interface {
	List() ([]*domain.Parameter, error)
	Get(code string) (*domain.Parameter, error)
	GetBatch(code ...string) ([]*domain.Parameter, error)
	Create(param *ParameterInfo) (*domain.Parameter, error)
	Update(param *ParameterInfo) (*domain.Parameter, error)
	Delete(code string) error
}
