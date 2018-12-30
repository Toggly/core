package mongo_test

import (
	"testing"

	"github.com/Toggly/core/storage"
	asserts "github.com/stretchr/testify/assert"
)

func TestMongoEnvironment(t *testing.T) {
	assert := asserts.New(t)

	beforeTest()

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

	afterTest()
}
