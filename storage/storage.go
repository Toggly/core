package storage

import (
	"errors"
	"fmt"

	"github.com/Toggly/core/domain"
)

// ErrUniqueIndex type
type ErrUniqueIndex struct {
	Type string
	Key  string
}

func (e ErrUniqueIndex) Error() string {
	return fmt.Sprintf("Unique index error: %s [%s]", e.Type, e.Key)
}

var (
	// ErrNotFound error
	ErrNotFound = errors.New("Not found")
	// ErrEntityRelationsBroken error
	ErrEntityRelationsBroken = errors.New("Entity relations broken")
)

// DataStorage defines storage interface
type DataStorage interface {
	Connect() error
	Projects(owner string) ProjectStorage
	Environments(owner, project string) EnvironmentStorage
	Groups(owner, project, env string) GroupStorage
	Parameters(owner, project, env, group string) ParameterStorage
}

// ProjectStorage defines projects storage interface
type ProjectStorage interface {
	List() ([]*domain.Project, error)
	Get(code string) (*domain.Project, error)
	Delete(code string) error
	Save(project *domain.Project) error
	Update(project *domain.Project) error
}

// EnvironmentStorage defines environments storage interface
type EnvironmentStorage interface {
	List() ([]*domain.Environment, error)
	Get(code string) (*domain.Environment, error)
	Delete(code string) error
	Save(env *domain.Environment) error
	Update(env *domain.Environment) error
}

// GroupStorage defines groups storage interface
type GroupStorage interface {
	List() ([]*domain.Group, error)
	Get(code string) (*domain.Group, error)
	Delete(code string) error
	Save(grp *domain.Group) error
	Update(grp *domain.Group) error
}

// ParameterStorage defines groups storage interface
type ParameterStorage interface {
	List() ([]*domain.Parameter, error)
	Get(code string) (*domain.Parameter, error)
	GetBatch(code ...string) ([]*domain.Parameter, error)
	Delete(code string) error
	Save(grp *domain.Parameter) error
	Update(grp *domain.Parameter) error
}
