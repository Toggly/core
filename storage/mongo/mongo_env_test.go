package mongo_test

import (
	"testing"

	"github.com/Toggly/core/domain"
	"github.com/Toggly/core/storage"
	"github.com/Toggly/core/util"
	asserts "github.com/stretchr/testify/assert"
)

func TestMongoEnvironment(t *testing.T) {
	assert := asserts.New(t)

	beforeTest()

	var err error

	db := getDB().Environments("ow1", "proj1")

	t.Run("get not found", func(t *testing.T) {
		env, err := db.Get("env1")
		assert.Nil(env)
		assert.Equal(storage.ErrNotFound, err)
	})

	t.Run("list empty", func(t *testing.T) {
		list, err := db.List()
		assert.Nil(err)
		assert.NotNil(list)
		assert.Len(list, 0)
	})

	t.Run("create", func(t *testing.T) {
		e := &domain.Environment{
			Owner:       "ow1",
			Project:     "proj1",
			Code:        "env1",
			Description: "Description 1",
			Protected:   true,
			RegDate:     util.Now(),
		}

		t.Run("wrong owner", func(t *testing.T) {
			e.Owner = "ow2"
			err = db.Save(e)
			assert.Equal(storage.ErrEntityRelationsBroken, err)
			e.Owner = "ow1"
		})

		t.Run("wrong project", func(t *testing.T) {
			e.Project = "proj2"
			err = db.Save(e)
			assert.Equal(storage.ErrEntityRelationsBroken, err)
			e.Project = "proj1"
		})

		t.Run("ok", func(t *testing.T) {
			err = db.Save(e)
			assert.Nil(err)

			env, err := db.Get("env1")
			assert.Nil(err)
			assert.Equal(e.Code, env.Code)
			assert.Equal(e.Description, env.Description)
			assert.Equal(e.RegDate, env.RegDate)
			assert.Equal(e.Owner, env.Owner)
			assert.Equal(e.Project, env.Project)
			assert.True(env.Protected)
		})

		t.Run("unique index error", func(t *testing.T) {
			err = db.Save(e)
			assert.NotNil(err)
			_, ok := err.(*storage.ErrUniqueIndex)
			assert.True(ok)
		})
	})

	t.Run("list one item", func(t *testing.T) {
		list, err := db.List()
		assert.Nil(err)
		assert.NotNil(list)
		assert.Len(list, 1)
	})

	t.Run("empty list for owner 2", func(t *testing.T) {
		list, err := getDB().Environments("ow2", "proj1").List()
		assert.Nil(err)
		assert.NotNil(list)
		assert.Len(list, 0)
	})

	t.Run("empty list for project 2", func(t *testing.T) {
		list, err := getDB().Environments("ow1", "proj2").List()
		assert.Nil(err)
		assert.NotNil(list)
		assert.Len(list, 0)
	})

	t.Run("update", func(t *testing.T) {
		e := &domain.Environment{
			Owner:       "ow1",
			Project:     "proj1",
			Code:        "env1",
			Description: "Description 2",
			Protected:   true,
			RegDate:     util.Now(),
		}

		t.Run("not found", func(t *testing.T) {
			e.Code = "env3"
			err = db.Update(e)
			assert.Equal(storage.ErrNotFound, err)
			e.Code = "env1"
		})

		t.Run("wrong owner", func(t *testing.T) {
			e.Owner = "ow2"
			err = db.Update(e)
			assert.Equal(storage.ErrEntityRelationsBroken, err)
			e.Owner = "ow1"
		})

		t.Run("wrong project", func(t *testing.T) {
			e.Project = "proj2"
			err = db.Update(e)
			assert.Equal(storage.ErrEntityRelationsBroken, err)
			e.Project = "proj1"
		})

		t.Run("ok", func(t *testing.T) {
			err = db.Update(e)
			assert.Nil(err)

			env1, err := db.Get("env1")
			assert.Nil(err)
			assert.Equal(e.Code, env1.Code)
			assert.Equal(e.Description, env1.Description)
			assert.Equal(e.RegDate, env1.RegDate)
			assert.Equal(e.Owner, env1.Owner)
			assert.True(env1.Protected)
		})
	})

	t.Run("delete", func(t *testing.T) {
		db.Delete("env1")
		e, err := db.Get("env1")
		assert.Nil(e)
		assert.Equal(storage.ErrNotFound, err)
	})

	afterTest()
}
