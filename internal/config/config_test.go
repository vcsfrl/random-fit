package config

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/core"
	"go.starlark.net/starlark"
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
	var collDict = starlark.NewDict(1)
	var metaDict = starlark.NewDict(1)
	var setDict = starlark.NewDict(1)
	var setList = starlark.NewList([]starlark.Value{setDict})

	_ = collDict.SetKey(starlark.String("sets"), setList)

	_ = metaDict.SetKey(starlark.String("id"), starlark.String("test-coll-id"))
	_ = metaDict.SetKey(starlark.String("name"), starlark.String("test-coll-name"))
	_ = metaDict.SetKey(starlark.String("description"), starlark.String("test collection description"))

	_ = collDict.SetKey(starlark.String("metadata"), metaDict)

	_ = setDict.SetKey(starlark.String("metadata"), metaDict)
	_ = setDict.SetKey(starlark.String("id"), starlark.String("test-set-id"))
	_ = setDict.SetKey(starlark.String("name"), starlark.String("test-set-name"))
	_ = setDict.SetKey(starlark.String("description"), starlark.String("test set description"))

	collJson := []byte(collDict.String())
	coll := core.Collection{}

	_ = json.Unmarshal(collJson, &coll)

	fmt.Printf("Collection: %+v\nSets: %+v\n", coll, coll.Sets)
	for _, set := range coll.Sets {
		fmt.Printf("Set: %+v\n", set)
	}
	fmt.Printf("Collection json: %v\n", collJson)
}
