package build

//import (
//	"github.com/stretchr/testify/suite"
//	"github.com/vcsfrl/random-fit/internal/combination"
//	"testing"
//)
//
//func TestBuildSuite(t *testing.T) {
//	suite.Run(t, new(BuildSuite))
//}
//
//type BuildSuite struct {
//	suite.Suite
//}
//
//func (suite *BuildSuite) TestBuild() {
//	definition := &combination.CombinationDefinition{
//		ID:         "test",
//		StarScript: "./testdata/star_script.star",
//	}
//
//	builder, err := NewBuilder(definition)
//	suite.NoError(err)
//
//	combination, err := builder.Build()
//
//	suite.NoError(err)
//	suite.NotNil(combination)
//	suite.Equal(definition, combination.Definition)
//	suite.NotEmpty(combination.UUID)
//	suite.NotNil(combination.Data)
//}
