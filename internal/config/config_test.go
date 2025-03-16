package config

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/core"
	"go.starlark.net/starlark"
	"log"
	"testing"
	"time"
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
