package build

import (
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/model"
	"testing"
)

func TestBuildSuite(t *testing.T) {
	suite.Run(t, new(BuildSuite))
}

type BuildSuite struct {
	suite.Suite
}

func (suite *BuildSuite) TestBuild() {
	definition := &model.CombinationDefinition{
		ID:              "test",
		StarScript:      "./testdata/star_script.star",
		GoTemplate:      "./testdata/go_template.gohtml",
		OutputExtension: "md",
		OutputMimeType:  "text/markdown",
	}

	builder, err := NewBuilder(definition)
	suite.NoError(err)

	combination, err := builder.Build()

	suite.NoError(err)
	suite.NotNil(combination)
	suite.Equal(definition, combination.Definition)
	suite.NotEmpty(combination.UUID)
	suite.NotNil(combination.Data)
}
