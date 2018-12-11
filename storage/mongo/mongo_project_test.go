package mongo_test

import (
	"testing"

	"github.com/Toggly/core/storage"

	"github.com/Toggly/core/domain"
	"github.com/Toggly/core/util"
	asserts "github.com/stretchr/testify/assert"
)

func TestMongoProject(t *testing.T) {
	assert := asserts.New(t)

	beforeTest()

	var err error

	db := getDB().ForOwner("ow1").Projects()

	t.Run("get not found", func(t *testing.T) {
		proj, err := db.Get("proj1")
		assert.Nil(proj)
		assert.Equal(storage.ErrNotFound, err)
	})

	t.Run("list empty", func(t *testing.T) {
		list, err := db.List()
		assert.Nil(err)
		assert.NotNil(list)
		assert.Len(list, 0)
	})

	t.Run("create", func(t *testing.T) {
		p := &domain.Project{
			Code:        "proj1",
			Description: "Description 1",
			Owner:       "ow1",
			Status:      domain.ProjectStatusActive,
			RegDate:     util.Now(),
		}

		t.Run("wrong owner", func(t *testing.T) {
			p.Owner = "ow2"
			err = db.Save(p)
			assert.Equal(storage.ErrEntityRelationsBroken, err)
		})

		t.Run("ok", func(t *testing.T) {
			p.Owner = "ow1"
			err = db.Save(p)
			assert.Nil(err)

			proj, err := db.Get("proj1")
			assert.Nil(err)
			assert.Equal(p.Code, proj.Code)
			assert.Equal(p.Description, proj.Description)
			assert.Equal(p.RegDate, proj.RegDate)
			assert.Equal(p.Owner, proj.Owner)
			assert.Equal(p.Status, proj.Status)
		})
	})

	t.Run("list one item", func(t *testing.T) {
		list, err := db.List()
		assert.Nil(err)
		assert.NotNil(list)
		assert.Len(list, 1)
	})

	t.Run("update", func(t *testing.T) {
		p := &domain.Project{
			Code:        "proj1",
			Description: "Description 2",
			Owner:       "ow1",
			Status:      domain.ProjectStatusDisabled,
			RegDate:     util.Now(),
		}

		t.Run("wrong owner", func(t *testing.T) {
			p.Owner = "ow2"
			err = db.Update(p)
			assert.Equal(storage.ErrEntityRelationsBroken, err)
		})

		t.Run("ok", func(t *testing.T) {
			p.Owner = "ow1"
			err = db.Update(p)
			assert.Nil(err)

			proj, err := db.Get("proj1")
			assert.Nil(err)
			assert.Equal(p.Code, proj.Code)
			assert.Equal(p.Description, proj.Description)
			assert.Equal(p.RegDate, proj.RegDate)
			assert.Equal(p.Owner, proj.Owner)
			assert.Equal(p.Status, proj.Status)
		})
	})

	t.Run("delete", func(t *testing.T) {
		db.Delete("proj1")
		proj, err := db.Get("proj1")
		assert.Nil(proj)
		assert.Equal(storage.ErrNotFound, err)
	})

	afterTest()
}
