package build

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestBuildSuite(t *testing.T) {
	suite.Run(t, new(BuildSuite))
}

type BuildSuite struct {
	suite.Suite

	testFolder string
	builder    *StartCollectionBuilder
}

func (suite *BuildSuite) SetupTest() {
	suite.testFolder = "testdata/"
	builder, err := NewStartCollectionBuilder("testdata/collection.star")
	suite.NoError(err)
	suite.builder = builder
}

func (suite *BuildSuite) TestFromScript() {

	collection, err := suite.builder.Build()
	suite.NoError(err)

	spew.Dump(collection)
}
