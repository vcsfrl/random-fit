package build

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func TestBuildSuite(t *testing.T) {
	suite.Run(t, new(BuildSuite))
}

type BuildSuite struct {
	suite.Suite

	builder  *StartCollectionBuilder
	id       int
	testRand int
}

func (suite *BuildSuite) SetupTest() {
	var err error
	suite.builder, err = NewStartCollectionBuilder("testdata/collection.star")
	suite.NoError(err)

	suite.id = 0
	suite.builder.uuidFunc = func() (string, error) {
		suite.id++
		return fmt.Sprintf("00000000-0000-0000-0000-%012d", suite.id), nil
	}

	suite.builder.nowFunc = func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	suite.builder.randomIntFunc = func(min int, max int) int {
		suite.testRand++
		return suite.testRand
	}

	suite.NoError(err)
}

func (suite *BuildSuite) TestFromScript() {
	collection, err := suite.builder.Build()
	suite.NoError(err)

	jsonData, err := json.MarshalIndent(collection, "", "  ")
	suite.NoError(err)

	suite.Equal(buildFromScriptResult1, string(jsonData))

	collection, err = suite.builder.Build()
	suite.NoError(err)
	jsonData, err = json.MarshalIndent(collection, "", "  ")
	suite.NoError(err)

	suite.Equal(buildFromScriptResult2, string(jsonData))
}
