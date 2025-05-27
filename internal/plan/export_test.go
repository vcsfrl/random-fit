package plan

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/internal/combination"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExportSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExportSuite))
}

type ExportSuite struct {
	suite.Suite

	testFolder         string
	combinationFolder  string
	storageFolder      string
	planDefinition     *Definition
	combinationBuilder combination.Builder
	planBuilder        *Builder
	id                 int
	testRand           uint
}

func (suite *ExportSuite) SetupTest() {
	suite.testFolder = filepath.Join("..", "..", "data", "test", uuid.New().String())

	// Create the test folder
	err := os.MkdirAll(suite.testFolder, 0755)
	suite.Require().NoError(err)

	suite.combinationFolder = filepath.Join(suite.testFolder, "combination")
	err = os.MkdirAll(suite.combinationFolder, 0755)
	suite.Require().NoError(err)

	suite.storageFolder = filepath.Join(suite.combinationFolder, "storage")
	err = os.MkdirAll(suite.storageFolder, 0755)
	suite.Require().NoError(err)

	// Create a test plan definition
	suite.planDefinition = &Definition{
		ID:      "test-definition",
		Details: "Test definition",
		Users:   []string{"user-1"},
		UserData: UserData{
			ContainerName:            []string{"GroupCombination-Container", "_date"},
			RecurrentGroupNamePrefix: "Recurrent-GroupCombination",
			RecurrentGroups:          4,
			NrOfGroupCombinations:    3,
		},
	}

	definition, err := combination.NewCombinationDefinition("./testdata/star_script.star")
	suite.Require().NoError(err)
	suite.combinationBuilder, err = combination.NewStarBuilder(definition)
	suite.Require().NoError(err)

	suite.id = 0

	definition.UUIDModule.SetUUIDFunc(func() (string, error) {
		suite.id++

		return fmt.Sprintf("00000000-0000-0000-0000-%012d", suite.id), nil
	})

	suite.testRand = 0

	definition.RandomModule.SetUintFunc(func(_ uint, _ uint) (uint, error) {
		suite.testRand++

		return suite.testRand, nil
	})

	suite.planBuilder = NewBuilder(suite.planDefinition, suite.combinationBuilder)
	suite.planBuilder.Now = func() time.Time {
		return time.Date(2010, 1, 2, 3, 4, 5, 6, time.UTC)
	}
}

func (suite *ExportSuite) TearDownTest() {
	// Remove the test folder
	err := os.RemoveAll(suite.testFolder)
	suite.Require().NoError(err)
}

func (suite *ExportSuite) TestExport() {
	plan, err := suite.planBuilder.Build()
	suite.Require().NoError(err)
	suite.NotNil(plan)

	exporter := NewExporter(suite.combinationFolder, suite.storageFolder)
	err = exporter.Export(plan)
	suite.Require().NoError(err)

	// Check if the user folder exists
	userFolder := filepath.Join(suite.combinationFolder, "user-1")
	suite.True(suite.fileExists(userFolder))

	// Check if the group folder exists
	groupContainer := filepath.Join(userFolder, "GroupCombination-Container", "2010-01-02-03-04")
	suite.True(suite.fileExists(groupContainer))

	for i := 1; i <= 4; i++ {
		groupFolder := filepath.Join(groupContainer, fmt.Sprintf("Recurrent-GroupCombination-%d", i))
		suite.True(suite.fileExists(groupFolder))

		extensions := []string{"json", "md"}

		// Check if the group combinations exist
		for j := 1; j <= 3; j++ {
			for _, ext := range extensions {
				groupCombination := filepath.Join(groupFolder, fmt.Sprintf("Lotto_Number_Picks_%d.%s", j, ext))
				exists, err := suite.fileExists(groupCombination)
				suite.Require().NoError(err)
				suite.True(exists, "File %s does not exist", groupCombination)

				// Check if the file is not empty
				fileInfo, err := os.Stat(groupCombination)
				suite.Require().NoError(err)
				suite.Positive(fileInfo.Size(), "File %s is empty", groupCombination)

				// Check if file contains a specific string
				file, err := os.ReadFile(groupCombination)
				suite.Require().NoError(err)
				suite.Contains(string(file), "Lotto Number Picks")
				suite.Contains(string(file), "Lotto Numbers for User 1")
				suite.Contains(string(file), "6/49 and Lucky Number")
			}
		}
	}
}

func (suite *ExportSuite) TestExportGenerate() {
	planGenerator := suite.planBuilder.Generate(context.Background())
	suite.NotNil(planGenerator)

	exporter := NewExporter(suite.combinationFolder, suite.storageFolder)
	err := exporter.ExportGenerator(context.Background(), planGenerator)
	suite.Require().NoError(err)

	// Check if the user folder exists
	userFolder := filepath.Join(suite.combinationFolder, "user-1")
	suite.True(suite.fileExists(userFolder))

	// Check if the group folder exists
	groupContainer := filepath.Join(userFolder, "GroupCombination-Container", "2010-01-02-03-04")
	suite.True(suite.fileExists(groupContainer))

	for i := 1; i <= 4; i++ {
		groupFolder := filepath.Join(groupContainer, fmt.Sprintf("Recurrent-GroupCombination-%d", i))
		suite.True(suite.fileExists(groupFolder))

		extensions := []string{"json", "md"}

		// Check if the group combinations exist
		for j := 1; j <= 3; j++ {
			for _, ext := range extensions {
				groupCombination := filepath.Join(groupFolder, fmt.Sprintf("Lotto_Number_Picks_%d.%s", j, ext))
				exists, err := suite.fileExists(groupCombination)
				suite.Require().NoError(err)
				suite.True(exists, "File %s does not exist")

				// Check if the file is not empty
				fileInfo, err := os.Stat(groupCombination)
				suite.Require().NoError(err)
				suite.Positive(fileInfo.Size(), "File %s is empty", groupCombination)

				// Check if file contains a specific string
				file, err := os.ReadFile(groupCombination)
				suite.Require().NoError(err)
				suite.Contains(string(file), "Lotto Number Picks")
				suite.Contains(string(file), "Lotto Numbers for User 1")
				suite.Contains(string(file), "6/49 and Lucky Number")
			}
		}
	}
}

func (suite *ExportSuite) TestExportNoDateInContainer() {
	suite.planDefinition.ContainerName = []string{"GroupCombination-Container"}
	plan, err := suite.planBuilder.Build()
	suite.Require().NoError(err)
	suite.NotNil(plan)

	exporter := NewExporter(suite.combinationFolder, suite.storageFolder)
	err = exporter.Export(plan)
	suite.Require().NoError(err)

	// Check if the user folder exists
	userFolder := filepath.Join(suite.combinationFolder, "user-1")
	suite.True(suite.fileExists(userFolder))

	// Check if the group folder exists.
	groupContainer := filepath.Join(userFolder, "GroupCombination-Container")
	suite.True(suite.fileExists(groupContainer))

	// Check that the date was not included in the container name.
	groupContainer = filepath.Join(userFolder, "GroupCombination-Container", "2010-01-02-03-04")
	suite.False(suite.fileExists(groupContainer))
}

func (suite *ExportSuite) TestExportObject() {
	plan, err := suite.planBuilder.Build()
	suite.Require().NoError(err)
	suite.NotNil(plan)

	exporter := NewExporter(suite.combinationFolder, suite.storageFolder)
	err = exporter.Export(plan)
	suite.Require().NoError(err)

	// Check if the user folder exists
	dataFile := filepath.Join(suite.storageFolder, plan.UUID.String()+".gob")
	suite.True(suite.fileExists(dataFile))

	// open the file
	file, err := os.Open(dataFile)
	suite.Require().NoError(err)

	savedPlan := &UserPlan{}

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(savedPlan)
	suite.Require().NoError(err)
	suite.Equal(plan.UUID, savedPlan.UUID)
	suite.Equal(plan.CreatedAt, savedPlan.CreatedAt)
	suite.Len(plan.UserGroups, len(savedPlan.UserGroups))
	suite.Len(plan.UserGroups["user-1"], len(savedPlan.UserGroups["user-1"]))
	suite.Equal(plan.UserGroups["user-1"][0].Details, savedPlan.UserGroups["user-1"][0].Details)
	suite.Equal(
		plan.UserGroups["user-1"][0].Combinations[0].UUID,
		savedPlan.UserGroups["user-1"][0].Combinations[0].UUID,
	)
	suite.Equal(
		plan.UserGroups["user-1"][0].Combinations[0].Details,
		savedPlan.UserGroups["user-1"][0].Combinations[0].Details,
	)
	suite.Equal(
		plan.UserGroups["user-1"][0].Combinations[0].DefinitionID,
		savedPlan.UserGroups["user-1"][0].Combinations[0].DefinitionID,
	)
	suite.Equal(
		plan.UserGroups["user-1"][0].Combinations[0].Data,
		savedPlan.UserGroups["user-1"][0].Combinations[0].Data,
	)
	suite.Equal(
		plan.UserGroups["user-1"][0].Combinations[0].CreatedAt.Format(time.DateTime),
		savedPlan.UserGroups["user-1"][0].Combinations[0].CreatedAt.Format(time.DateTime),
	)

	err = file.Close()
	suite.Require().NoError(err)
}

func (suite *ExportSuite) TestExportObjectInFolder() {
	planGenerator := suite.planBuilder.Generate(context.Background())
	suite.NotNil(planGenerator)

	exporter := NewExporter(suite.combinationFolder, suite.storageFolder)
	err := exporter.ExportGenerator(context.Background(), planGenerator)
	suite.Require().NoError(err)

	// Check if the user folder exists
	userFolder := filepath.Join(suite.combinationFolder, "user-1")
	suite.True(suite.fileExists(userFolder))

	// get files from storage folder
	files, err := os.ReadDir(suite.storageFolder)
	suite.Require().NoError(err)
	suite.NotEmpty(files)

	suite.Len(files, suite.planDefinition.RecurrentGroups*suite.planDefinition.NrOfGroupCombinations)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		suite.Contains(file.Name(), suite.planDefinition.Users[0])
		suite.Contains(file.Name(), ".gob")

		// get file content
		dataFile := filepath.Join(suite.storageFolder, file.Name())
		suite.True(suite.fileExists(dataFile))
		// open the file
		data, err := os.ReadFile(dataFile)
		suite.Require().NoError(err)
		suite.NotEmpty(data)
		suite.Contains(string(data), "Lotto Number Picks")
		suite.Contains(string(data), "Lotto Numbers for User 1")
		suite.Contains(string(data), "Lucky Number")
	}
}

func (suite *ExportSuite) TestExportObjectInFolderCancelContext() {
	planGenerator := suite.planBuilder.Generate(context.Background())
	suite.NotNil(planGenerator)

	exporter := NewExporter(suite.combinationFolder, suite.storageFolder)
	background := context.Background()
	ctx, cancel := context.WithCancel(background)
	cancel()

	err := exporter.ExportGenerator(ctx, planGenerator)
	suite.Require().Error(err)
	suite.Equal(ErrExportTerminated, err)
}

func (suite *ExportSuite) fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, fmt.Errorf("error checking file existence: %w", err)
}
