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

	testFolder string
	builder    *StartCollectionBuilder
}

func (suite *BuildSuite) SetupTest() {
	suite.testFolder = "testdata/"
	id := 0

	var err error

	suite.builder, err = NewStartCollectionBuilder("testdata/collection.star")
	suite.NoError(err)

	suite.builder.uuidFunc = func() (string, error) {
		id++
		return fmt.Sprintf("00000000-0000-0000-0000-%012d", id), nil
	}

	suite.builder.nowFunc = func() time.Time {
		return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	suite.NoError(err)
}

func (suite *BuildSuite) TestFromScript() {
	collection, err := suite.builder.Build()
	suite.NoError(err)

	jsonData, err := json.MarshalIndent(collection, "", "  ")
	suite.NoError(err)

	suite.Equal(string(jsonData), buildFromScriptResult1)

	collection, err = suite.builder.Build()
	suite.NoError(err)
	jsonData, err = json.MarshalIndent(collection, "", "  ")
	suite.NoError(err)

	suite.Equal(string(jsonData), buildFromScriptResult2)
}
