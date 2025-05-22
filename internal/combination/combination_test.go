package combination_test

import (
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/combination"
	"testing"
)

func TestCombinationDataSuite(t *testing.T) {
	suite.Run(t, new(CombinationDataSuite))
}

type CombinationDataSuite struct {
	suite.Suite
}

func (suite *CombinationDataSuite) TestUnmarshall() {
	jsonData := []byte(`{
		"Extension": "json",
		"MimeType": "application/json",
		"Type": "json",
		"Data": "{\"numbers\": [1, 2, 3, 4, 5, 6]}"
	}`)

	var data combination.Data
	err := data.UnmarshalJSON(jsonData)
	suite.NoError(err)

	suite.Equal("json", data.Extension)
	suite.Equal("application/json", data.MimeType)
	suite.Equal(combination.DataTypeJson, data.Type)
	suite.Equal(`{"numbers": [1, 2, 3, 4, 5, 6]}`, data.Data.String())
}
