package config

import (
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

func (suite *ConfigSuite) TestFromFile() {

	// When
	elements := ElementsFromFolder(suite.testFolder)

	// Then
	suite.NotNil(elements)
	suite.Equal(0, len(elements))
}
