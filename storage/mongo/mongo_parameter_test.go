package mongo_test

import (
	"testing"

	asserts "github.com/stretchr/testify/assert"
)

func TestMongoParameter(t *testing.T) {
	t.Skip("Not implemented")
	assert := asserts.New(t)

	beforeTest()

	var err error

	db := getDB().Parameters("ow1", "proj1", "env1", "grp1")

	_ = assert
	_ = err
	_ = db

	afterTest()
}
