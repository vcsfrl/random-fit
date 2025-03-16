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

	suite.builder = &StartCollectionBuilder{
		starFile: "testdata/collection.star",
		uuidFunc: func() (string, error) {
			id++
			return fmt.Sprintf("00000000-0000-0000-0000-%012d", id), nil
		},
		nowFunc: func() time.Time {
			return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		},
	}

	err := suite.builder.start()
	suite.NoError(err)
}

func (suite *BuildSuite) TestFromScript() {

	collection, err := suite.builder.Build()
	suite.NoError(err)

	jsonData, err := json.MarshalIndent(collection, "", "  ")
	suite.NoError(err)

	suite.Equal(string(jsonData), `{
  "Metadata": {
    "ID": "collection-00000000-0000-0000-0000-000000000001",
    "Name": "Lotto number picks",
    "Description": "Users monthly Lotto Number picks",
    "Date": "2021-01-01T00:00:00Z"
  },
  "Sets": [],
  "Collections": [
    {
      "Metadata": {
        "ID": "collection-00000000-0000-0000-0000-000000000002",
        "Name": "Lotto Numbers fot User 1",
        "Description": "User 1 monthly Lotto Number picks",
        "Date": "2021-01-01T00:00:00Z"
      },
      "Sets": [
        {
          "Metadata": {
            "ID": "set-00000000-0000-0000-0000-000000000003",
            "Name": "Lotto Numbers fot User 1",
            "Description": "User 1 monthly Lotto Number picks",
            "Date": "2021-01-01T00:00:00Z"
          },
          "Elements": [
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000004",
                "Name": "Numbers",
                "Description": "6 numbers out of 49",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                1,
                2,
                3,
                4,
                5,
                6
              ]
            },
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000005",
                "Name": "Lucky Number",
                "Description": "Lucky Number for 6/49 draw",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                24500
              ]
            }
          ]
        },
        {
          "Metadata": {
            "ID": "set-00000000-0000-0000-0000-000000000006",
            "Name": "Lotto Numbers fot User 1",
            "Description": "User 1 monthly Lotto Number picks",
            "Date": "2021-01-01T00:00:00Z"
          },
          "Elements": [
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000007",
                "Name": "Numbers",
                "Description": "6 numbers out of 49",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                1,
                2,
                3,
                4,
                5,
                6
              ]
            },
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000008",
                "Name": "Lucky Number",
                "Description": "Lucky Number for 6/49 draw",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                24500
              ]
            }
          ]
        },
        {
          "Metadata": {
            "ID": "set-00000000-0000-0000-0000-000000000009",
            "Name": "Lotto Numbers fot User 1",
            "Description": "User 1 monthly Lotto Number picks",
            "Date": "2021-01-01T00:00:00Z"
          },
          "Elements": [
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000010",
                "Name": "Numbers",
                "Description": "6 numbers out of 49",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                1,
                2,
                3,
                4,
                5,
                6
              ]
            },
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000011",
                "Name": "Lucky Number",
                "Description": "Lucky Number for 6/49 draw",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                24500
              ]
            }
          ]
        }
      ],
      "Collections": null
    },
    {
      "Metadata": {
        "ID": "collection-00000000-0000-0000-0000-000000000012",
        "Name": "Lotto Numbers fot User 2",
        "Description": "User 2 monthly Lotto Number picks",
        "Date": "2021-01-01T00:00:00Z"
      },
      "Sets": [
        {
          "Metadata": {
            "ID": "set-00000000-0000-0000-0000-000000000013",
            "Name": "Lotto Numbers fot User 2",
            "Description": "User 2 monthly Lotto Number picks",
            "Date": "2021-01-01T00:00:00Z"
          },
          "Elements": [
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000014",
                "Name": "Numbers",
                "Description": "6 numbers out of 49",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                1,
                2,
                3,
                4,
                5,
                6
              ]
            },
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000015",
                "Name": "Lucky Number",
                "Description": "Lucky Number for 6/49 draw",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                24500
              ]
            }
          ]
        },
        {
          "Metadata": {
            "ID": "set-00000000-0000-0000-0000-000000000016",
            "Name": "Lotto Numbers fot User 2",
            "Description": "User 2 monthly Lotto Number picks",
            "Date": "2021-01-01T00:00:00Z"
          },
          "Elements": [
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000017",
                "Name": "Numbers",
                "Description": "6 numbers out of 49",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                1,
                2,
                3,
                4,
                5,
                6
              ]
            },
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000018",
                "Name": "Lucky Number",
                "Description": "Lucky Number for 6/49 draw",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                24500
              ]
            }
          ]
        },
        {
          "Metadata": {
            "ID": "set-00000000-0000-0000-0000-000000000019",
            "Name": "Lotto Numbers fot User 2",
            "Description": "User 2 monthly Lotto Number picks",
            "Date": "2021-01-01T00:00:00Z"
          },
          "Elements": [
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000020",
                "Name": "Numbers",
                "Description": "6 numbers out of 49",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                1,
                2,
                3,
                4,
                5,
                6
              ]
            },
            {
              "Metadata": {
                "ID": "element-00000000-0000-0000-0000-000000000021",
                "Name": "Lucky Number",
                "Description": "Lucky Number for 6/49 draw",
                "Date": "2021-01-01T00:00:00Z"
              },
              "Values": [
                24500
              ]
            }
          ]
        }
      ],
      "Collections": null
    }
  ]
}`,
	)
}
