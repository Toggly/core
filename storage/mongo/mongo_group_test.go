package mongo_test

import (
	"testing"

	asserts "github.com/stretchr/testify/assert"
)

func TestMongoGroup(t *testing.T) {
	t.Skip("Not implemented")
	assert := asserts.New(t)

	beforeTest()

	var err error

	db := getDB().Groups("ow1", "proj1", "env1")

	_ = assert
	_ = err
	_ = db

	afterTest()
}
