package combination

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCombinationDataSuite(t *testing.T) {
	suite.Run(t, new(CombinationDataSuite))
}

type CombinationDataSuite struct {
	suite.Suite
}

func (suite *CombinationDataSuite) TestUnmarshall() {
	data := []byte(`{
		"Extension": "json",
		"MimeType": "application/json",
		"Type": "json",
		"Data": "{\"numbers\": [1, 2, 3, 4, 5, 6]}"
	}`)

	var d Data
	err := d.UnmarshalJSON(data)
	suite.NoError(err)

	suite.Equal("json", d.Extension)
	suite.Equal("application/json", d.MimeType)
	suite.Equal(DataTypeJson, d.Type)
	suite.Equal(`{"numbers": [1, 2, 3, 4, 5, 6]}`, d.Data.String())
}
