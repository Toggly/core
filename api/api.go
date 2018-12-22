package api

import (
	"errors"
	"fmt"

	"github.com/Toggly/core/domain"
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

// TogglyAPI interface
type TogglyAPI interface {
	ForOwner(owner string) OwnerAPI
}

// OwnerAPI interface
type OwnerAPI interface {
	Projects() ProjectAPI
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
	For(code string) ForProjectAPI
}

// ForProjectAPI interface
type ForProjectAPI interface {
	Environments() EnvironmentAPI
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
}

// ForEnvironmentAPI interface
type ForEnvironmentAPI interface {
	Groups() GroupAPI
	Parameters() ParameterAPI
}

// GroupInfo type
type GroupInfo struct{}

// GroupAPI interface
type GroupAPI interface {
	List() ([]*domain.Group, error)
	Get(code string) (domain.Group, error)
	Create(code string) (domain.Group, error)
	For(code string) ForGroupAPI
}

// ForGroupAPI interface
type ForGroupAPI interface {
	Parameters() ParameterAPI
}

// ParameterInfo type
type ParameterInfo struct{}

// ParameterAPI interface
type ParameterAPI interface {
	List() ([]*domain.Parameter, error)
	Get(code string) ([]*domain.Parameter, error)
	GetBatch(code ...string) ([]*domain.Parameter, error)
	Create(param *ParameterInfo) (domain.Parameter, error)
	Update(param *ParameterInfo) (domain.Parameter, error)
	Delete(code string) error
}
