package cmd_test

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"github.com/vcsfrl/random-fit/cmd"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestCommands(t *testing.T) {
	suite.Run(t, new(CommandsSuite))
}

type CommandsSuite struct {
	suite.Suite

	buffer  *bytes.Buffer
	command *cobra.Command

	testFolder    string
	codeGenFolder string
}

func (suite *CommandsSuite) SetupTest() {
	// Capture output
	var err error

	suite.buffer = new(bytes.Buffer)
	suite.command, err = cmd.NewCommand()
	suite.NoError(err)

	suite.command.SetOut(suite.buffer)
	suite.command.SetErr(suite.buffer)

	suite.testFolder = filepath.Join("..", "data", "test", uuid.New().String())

	// Create the test folder
	err = os.MkdirAll(suite.testFolder, 0755)
	suite.NoError(err)

	// Set the environment variable
	err = os.Setenv("RF_BASE_FOLDER", suite.testFolder)
	suite.NoError(err)
	err = os.Setenv("RF_DATA_FOLDER", suite.testFolder)
	suite.NoError(err)
	err = os.Setenv("EDITOR", "")
	suite.NoError(err)
}

func (suite *CommandsSuite) TearDownTest() {
	// Remove the test folder
	err := os.RemoveAll(suite.testFolder)
	suite.NoError(err)
}

func (suite *CommandsSuite) TestSubcommands() {
	subcommands := suite.command.Commands()
	var expectedSubcommandNames = []string{"definition", "code", "generate", "combination", "new", "edit", "delete", "plan", "new", "edit", "delete", "generate", "combination"}
	var subcommandNames []string

	for _, cmd := range subcommands {
		subcommandNames = append(subcommandNames, cmd.Name())
		for _, subCmd := range cmd.Commands() {
			subcommandNames = append(subcommandNames, subCmd.Name())
			for _, subSubCmd := range subCmd.Commands() {
				subcommandNames = append(subcommandNames, subSubCmd.Name())
			}
		}
	}

	slices.Sort(subcommandNames)
	slices.Sort(expectedSubcommandNames)

	suite.Equal(expectedSubcommandNames, subcommandNames)
}

func (suite *CommandsSuite) TestGenerateCode() {
	// Create the code generation folder
	suite.codeGenFolder = filepath.Join(suite.testFolder, "internal", "service")
	err := cmd.CreateFolder(suite.codeGenFolder)
	suite.NoError(err)

	// Create the internal folder to copy file that is used as a source
	codeFolder := filepath.Join(suite.testFolder, "internal", "combination", "template")
	err = cmd.CreateFolder(codeFolder)
	suite.NoError(err)

	// copy file template file to the code generation folder
	sourceFile := filepath.Join("..", "internal", "combination", "template", "script.star")
	destFile := filepath.Join(codeFolder, "script.star")
	// copy the file
	templateData, err := os.ReadFile(sourceFile)
	suite.NoError(err)
	err = os.WriteFile(destFile, templateData, 0644)
	suite.NoError(err)

	suite.command.SetArgs([]string{"code", "generate"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check generated file content
	generatedFile := filepath.Join(suite.codeGenFolder, "combination_definition_template.go")
	genData, err := os.ReadFile(generatedFile)
	suite.NoError(err)
	suite.NotEmpty(genData)
	suite.Contains(string(genData), string(templateData))
}

func (suite *CommandsSuite) TestDefinitionCombination_New() {
	suite.command.SetArgs([]string{"definition", "combination", "new"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, cmd.MsgNameMissing)

	suite.command.SetArgs([]string{"definition", "combination", "new", "--name", "test1"})
	err1 := suite.command.Execute()
	suite.NoError(err1)

	//// Check output
	scriptName := filepath.Join(suite.testFolder, "definition", "test1.star")
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgCreate+" "+cmd.MsgCombinationDefinition+" test1")
	suite.Contains(output, cmd.MsgDone+" "+cmd.MsgCreate+" "+cmd.MsgCombinationDefinition+" test1")
	suite.Contains(output, cmd.MsgEditScript+" "+scriptName)
	suite.Contains(output, cmd.ErrNoEnvEditor.Error())

	scriptData, err := os.ReadFile(scriptName)
	suite.NoError(err)
	suite.Contains(string(scriptData), "definition =")

	suite.command.SetArgs([]string{"definition", "combination", "new", "test2"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName = filepath.Join(suite.testFolder, "definition", "test2.star")
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgCreate+" "+cmd.MsgCombinationDefinition+" test2")
	suite.Contains(output, cmd.MsgDone+" "+cmd.MsgCreate+" "+cmd.MsgCombinationDefinition+" test2")
	suite.Contains(output, cmd.MsgEditScript+" "+scriptName)
	suite.Contains(output, cmd.ErrNoEnvEditor.Error())

	scriptData, err = os.ReadFile(scriptName)
	suite.NoError(err)
	suite.Contains(string(scriptData), "definition =")
}

func (suite *CommandsSuite) TestDefinitionCombination_Edit() {
	suite.command.SetArgs([]string{"definition", "combination", "edit"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, cmd.MsgNameMissing)

	suite.command.SetArgs([]string{"definition", "combination", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "definition", "test1.star")
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgEditScript+" "+scriptName)

	suite.Contains(output, cmd.ErrNoEnvEditor.Error())
}

func (suite *CommandsSuite) TestDefinitionCombination_Delete() {
	suite.command.SetArgs([]string{"definition", "combination", "delete"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, cmd.MsgNameMissing)

	suite.command.SetArgs([]string{"definition", "combination", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "definition", "test1.star")
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgEditScript+" "+scriptName)

	suite.command.SetArgs([]string{"definition", "combination", "delete", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)
	// Check output
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgDelete+" "+cmd.MsgCombinationDefinition+" test1")
	suite.Contains(output, cmd.MsgDone+" "+cmd.MsgDelete+" "+cmd.MsgCombinationDefinition+" test1")

	// check if the file is deleted
	_, err = os.Stat(scriptName)
	suite.True(os.IsNotExist(err), "File should be deleted")
}

func (suite *CommandsSuite) TestDefinitionCombination_List() {
	suite.command.SetArgs([]string{"definition", "combination"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, cmd.MsgCombinationDefinition+" "+cmd.MsgList)
	suite.Contains(output, cmd.MsgNoItemsFound)

	// Create a definition
	suite.command.SetArgs([]string{"definition", "combination", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)
	suite.command.SetArgs([]string{"definition", "combination", "new", "--name", "test2"})
	err = suite.command.Execute()
	suite.NoError(err)

	suite.command.SetArgs([]string{"definition", "combination"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	output = suite.buffer.String()
	suite.Contains(output, " - test1")
	suite.Contains(output, " - test2")
}

func (suite *CommandsSuite) TestDefinitionPlan_New() {
	suite.command.SetArgs([]string{"definition", "plan", "new"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, cmd.MsgNameMissing)

	suite.command.SetArgs([]string{"definition", "plan", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "plan", "test1.json")
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgCreate+" "+cmd.MsgPlanDefinition+" test1")
	suite.Contains(output, cmd.MsgDone+" "+cmd.MsgCreate+" "+cmd.MsgPlanDefinition+" test1")
	suite.Contains(output, cmd.MsgEditScript+" "+scriptName)
	suite.Contains(output, cmd.ErrNoEnvEditor.Error())

	scriptData, err := os.ReadFile(scriptName)
	suite.NoError(err)
	suite.Contains(string(scriptData), "RecurrentGroupNamePrefix")
}

func (suite *CommandsSuite) TestDefinitionPlan_List() {
	suite.command.SetArgs([]string{"definition", "plan"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, cmd.MsgPlanDefinition+" "+cmd.MsgList)
	suite.Contains(output, cmd.MsgNoItemsFound)

	// Create a definition
	suite.command.SetArgs([]string{"definition", "plan", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)
	suite.command.SetArgs([]string{"definition", "plan", "new", "--name", "test2"})
	err = suite.command.Execute()
	suite.NoError(err)

	suite.command.SetArgs([]string{"definition", "plan"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	output = suite.buffer.String()
	suite.Contains(output, " - test1")
	suite.Contains(output, " - test2")
}

func (suite *CommandsSuite) TestDefinitionPlan_Edit() {
	suite.command.SetArgs([]string{"definition", "plan", "edit"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, cmd.MsgNameMissing)

	suite.command.SetArgs([]string{"definition", "plan", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "plan", "test1.json")
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgEditScript+" "+scriptName)

	suite.Contains(output, cmd.ErrNoEnvEditor.Error())
}

func (suite *CommandsSuite) TestDefinitionPlan_Delete() {
	suite.command.SetArgs([]string{"definition", "plan", "delete"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, cmd.MsgNameMissing)

	suite.command.SetArgs([]string{"definition", "plan", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "plan", "test1.json")
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgEditScript+" "+scriptName)

	suite.command.SetArgs([]string{"definition", "plan", "delete", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgDelete+" "+cmd.MsgPlanDefinition+" test1")
	suite.Contains(output, cmd.MsgDone+" "+cmd.MsgDelete+" "+cmd.MsgPlanDefinition+" test1")
	// check if the file is deleted
	_, err = os.Stat(scriptName)
	suite.True(os.IsNotExist(err), "File should be deleted")
}

func (suite *CommandsSuite) TestGenerate_Combination() {
	suite.command.SetArgs([]string{"generate", "combination"})
	err := suite.command.Execute()
	suite.NoError(err)

	// check output
	output := suite.buffer.String()
	suite.Contains(output, cmd.MsgCombinationDefinitionNameMissing)

	suite.command.SetArgs([]string{"generate", "combination", "--combination", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// check output
	output = suite.buffer.String()
	suite.Contains(output, cmd.MsgPlanDefinitionNameMissing)

	// create combination definition
	suite.command.SetArgs([]string{"definition", "combination", "new", "--name", "combination1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// create plan definition
	suite.command.SetArgs([]string{"definition", "plan", "new", "--name", "plan1"})
	err = suite.command.Execute()
	suite.NoError(err)

	suite.command.SetArgs([]string{"generate", "combination", "--combination", "combination1", "--plan", "plan1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// get directories from test folder
	containerFolder := filepath.Join(suite.testFolder, "combination", "user1", "ContainerName")
	dirs, err := os.ReadDir(containerFolder)
	suite.NoError(err)
	suite.Len(dirs, 1)

	// check combination was created
	combinationMdFile := filepath.Join(containerFolder, dirs[0].Name(), "Group-1", "Sample_Combination_1.md")
	combinationJsonFile := filepath.Join(containerFolder, dirs[0].Name(), "Group-1", "Sample_Combination_1.json")

	_, err = os.Stat(combinationMdFile)
	suite.NoError(err, "File should exist")
	_, err = os.Stat(combinationJsonFile)
	suite.NoError(err, "File should exist")

	// check files content
	combinationMdData, err := os.ReadFile(combinationMdFile)
	suite.NoError(err)
	suite.Contains(string(combinationMdData), "Sample")

	combinationJsonData, err := os.ReadFile(combinationJsonFile)
	suite.NoError(err)
	suite.Contains(string(combinationJsonData), "Sample")

	// check objects saved in storage
	storageFolder := filepath.Join(suite.testFolder, "storage")

	// get files from storage folder
	storageFiles, err := os.ReadDir(storageFolder)
	suite.NoError(err)
	suite.Len(storageFiles, 1)
	// check files content
	storageFile := filepath.Join(storageFolder, storageFiles[0].Name())
	storageData, err := os.ReadFile(storageFile)
	suite.NoError(err)
	suite.Contains(string(storageData), "Sample")
}
