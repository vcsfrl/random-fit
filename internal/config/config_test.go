package config

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

type ConfigSuite struct {
	suite.Suite

	testFolder string
}

func (suite *ConfigSuite) SetupTest() {
	suite.testFolder = "testdata/"
}

func (suite *ConfigSuite) TestFromScript() {

	builder := NewStartCollectionBuilder("testdata/collection.star")
	builder.Start()
	collection := builder.Build()

	spew.Dump(collection)
}
