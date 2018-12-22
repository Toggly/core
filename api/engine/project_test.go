package engine_test

import (
	"testing"
	"time"

	"github.com/Toggly/core/domain"

	"github.com/Toggly/core/api"

	"github.com/Toggly/core/api/engine"
	asserts "github.com/stretchr/testify/assert"
)

func TestAPIProject(t *testing.T) {

	assert := asserts.New(t)
	e := engine.NewTogglyAPI(getDB(), logger)
	pApi := e.ForOwner("ow1").Projects()

	beforeTest()

	t.Run("get not found", func(t *testing.T) {
		proj, err := pApi.Get("proj1")
		assert.Nil(proj)
		assert.Equal(api.ErrProjectNotFound, err)
	})

	t.Run("list empty", func(t *testing.T) {
		list, err := pApi.List()
		assert.Nil(err)
		assert.Len(list, 0)
	})

	var regDate time.Time

	t.Run("bad request", func(t *testing.T) {
		tt := []*api.ProjectInfo{
			&api.ProjectInfo{},
			&api.ProjectInfo{Code: "p1", Status: "wrong"},
		}
		var err error
		for _, tc := range tt {
			_, err = pApi.Create(tc)
			_, ok := err.(*api.ErrBadRequest)
			assert.True(ok)
			_, err = pApi.Update(tc)
			_, ok = err.(*api.ErrBadRequest)
			assert.True(ok)
		}
	})

	t.Run("create", func(t *testing.T) {
		p := &api.ProjectInfo{
			Code:        "proj1",
			Description: "Project 1",
			Status:      domain.ProjectStatusActive,
		}
		proj, err := pApi.Create(p)
		assert.Nil(err)
		assert.NotNil(proj)
		assert.Equal(p.Code, proj.Code)
		assert.Equal(p.Description, proj.Description)
		assert.Equal(p.Status, proj.Status)
		assert.NotNil(proj.RegDate)
		regDate = proj.RegDate
		assert.Equal("ow1", proj.Owner)
	})

	t.Run("list one item", func(t *testing.T) {
		list, err := pApi.List()
		assert.Nil(err)
		assert.Len(list, 1)
	})

	t.Run("update", func(t *testing.T) {
		p := &api.ProjectInfo{
			Code:        "proj1",
			Description: "Project 2",
			Status:      domain.ProjectStatusDisabled,
		}
		proj, err := pApi.Update(p)
		assert.Nil(err)
		assert.NotNil(proj)
		assert.Equal(p.Code, proj.Code)
		assert.Equal(p.Description, proj.Description)
		assert.Equal(p.Status, proj.Status)
		assert.Equal("ow1", proj.Owner)
		assert.Equal(regDate, proj.RegDate)
	})

	t.Run("delete", func(t *testing.T) {
		err := pApi.Delete("proj1")
		assert.Nil(err)
		_, err = pApi.Get("proj1")
		assert.Equal(api.ErrProjectNotFound, err)
	})

	afterTest()

}
