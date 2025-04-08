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
	builder *Builder
}

func (suite *BuildSuite) TestBuild() {
	definition := &model.Definition{
		ID:         "test",
		StarScript: "testdata/star_script.star",
		GoTemplate: "./testdata/go_template.tmpl",
	}

	builder, err := NewBuilder(definition)
	suite.NoError(err)

	combination := builder.Build()

	suite.NotNil(combination)
	suite.Equal(definition.ID, combination.DefinitionId)
	suite.NotEmpty(combination.UUID)
	suite.Nil(combination.Data)
}
