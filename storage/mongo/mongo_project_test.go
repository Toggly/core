package mongo_test

import (
	"testing"

	"github.com/Toggly/core/domain"
	"github.com/Toggly/core/util"
	asserts "github.com/stretchr/testify/assert"
)

func TestMongoProject(t *testing.T) {
	assert := asserts.New(t)

	beforeTest()

	var err error

	db := getDB()

	p := &domain.Project{
		Code:        "proj1",
		Description: "Description 1",
		Owner:       "ow1",
		Status:      domain.ProjectStatusActive,
		RegDate:     util.Now(),
	}

	err = db.ForOwner("ow1").Projects().Save(p)

	_ = err

	proj, err := db.ForOwner("ow1").Projects().Get("proj1")
	assert.Nil(err)
	assert.Equal(p.Code, proj.Code)
	assert.Equal(p.Description, proj.Description)
	assert.Equal(p.RegDate, proj.RegDate)

	// db.ForOwner("ow1").Projects().List()

	// afterTest()
}
