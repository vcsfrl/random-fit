package cmd

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
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
	suite.command, err = NewCommand()
	suite.NoError(err)

	suite.command.SetOut(suite.buffer)
	suite.command.SetErr(suite.buffer)

	suite.testFolder = filepath.Join("..", "data", "test", uuid.New().String())

	// Create the test folder
	err = os.MkdirAll(suite.testFolder, 0755)
	suite.NoError(err)

	// Set the environment variable
	err = os.Setenv("RF_BASE_FOLDER", suite.testFolder)
	err = os.Setenv("RF_DATA_FOLDER", suite.testFolder)
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
	var expectedSubcommandNames = []string{"definition", "code", "generate", "combination", "new", "edit", "delete", "plan", "new", "edit", "delete"}
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
	suite.codeGenFolder = filepath.Join(suite.testFolder, "cmd", "internal")
	err := createFolder(suite.codeGenFolder)
	suite.NoError(err)

	// Create the internal folder to copy file that is used as a source
	codeFolder := filepath.Join(suite.testFolder, "internal", "combination", "template")
	err = createFolder(codeFolder)
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
	suite.Contains(output, msgNameMissing)

	suite.command.SetArgs([]string{"definition", "combination", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "definition", "test1.star")
	output = suite.buffer.String()
	suite.Contains(output, msgCreate+" "+msgCombinationDefinition+" test1")
	suite.Contains(output, msgDone+" "+msgCreate+" "+msgCombinationDefinition+" test1")
	suite.Contains(output, msgEditScript+" "+scriptName)
	suite.Contains(output, errNoEnvEditor.Error())
	scriptData, err := os.ReadFile(scriptName)
	suite.NoError(err)
	suite.Contains(string(scriptData), "definition =")

	suite.command.SetArgs([]string{"definition", "combination", "new", "test2"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName = filepath.Join(suite.testFolder, "definition", "test2.star")
	output = suite.buffer.String()
	suite.Contains(output, msgCreate+" "+msgCombinationDefinition+" test2")
	suite.Contains(output, msgDone+" "+msgCreate+" "+msgCombinationDefinition+" test2")
	suite.Contains(output, msgEditScript+" "+scriptName)
	suite.Contains(output, errNoEnvEditor.Error())
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
	suite.Contains(output, msgNameMissing)

	suite.command.SetArgs([]string{"definition", "combination", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "definition", "test1.star")
	output = suite.buffer.String()
	suite.Contains(output, msgEditScript+" "+scriptName)

	suite.Contains(output, errNoEnvEditor.Error())
}

func (suite *CommandsSuite) TestDefinitionCombination_Delete() {
	suite.command.SetArgs([]string{"definition", "combination", "delete"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, msgNameMissing)

	suite.command.SetArgs([]string{"definition", "combination", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "definition", "test1.star")
	output = suite.buffer.String()
	suite.Contains(output, msgEditScript+" "+scriptName)

	suite.command.SetArgs([]string{"definition", "combination", "delete", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)
	// Check output
	output = suite.buffer.String()
	suite.Contains(output, msgDelete+" "+msgCombinationDefinition+" test1")
	suite.Contains(output, msgRemoveScript+" "+scriptName)

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
	suite.Contains(output, msgCombinationDefinition+" "+msgList)
	suite.Contains(output, msgNoItemsFound)

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
	suite.Contains(output, msgNameMissing)

	suite.command.SetArgs([]string{"definition", "plan", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "plan", "test1.json")
	output = suite.buffer.String()
	suite.Contains(output, msgCreate+" "+msgPlanDefinition+" test1")
	suite.Contains(output, msgDone+" "+msgCreate+" "+msgPlanDefinition+" test1")
	suite.Contains(output, msgEditScript+" "+scriptName)
	suite.Contains(output, errNoEnvEditor.Error())
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
	suite.Contains(output, msgPlanDefinition+" "+msgList)
	suite.Contains(output, msgNoItemsFound)

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
	suite.Contains(output, msgNameMissing)

	suite.command.SetArgs([]string{"definition", "plan", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "plan", "test1.json")
	output = suite.buffer.String()
	suite.Contains(output, msgEditScript+" "+scriptName)

	suite.Contains(output, errNoEnvEditor.Error())
}

func (suite *CommandsSuite) TestDefinitionPlan_Delete() {
	suite.command.SetArgs([]string{"definition", "plan", "delete"})
	err := suite.command.Execute()
	suite.NoError(err)

	// Check output
	output := suite.buffer.String()
	suite.Contains(output, msgNameMissing)

	suite.command.SetArgs([]string{"definition", "plan", "new", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	scriptName := filepath.Join(suite.testFolder, "plan", "test1.json")
	output = suite.buffer.String()
	suite.Contains(output, msgEditScript+" "+scriptName)

	suite.command.SetArgs([]string{"definition", "plan", "delete", "--name", "test1"})
	err = suite.command.Execute()
	suite.NoError(err)

	// Check output
	output = suite.buffer.String()
	suite.Contains(output, msgDelete+" "+msgPlanDefinition+" test1")
	suite.Contains(output, msgRemoveScript+" "+scriptName)
	// check if the file is deleted
	_, err = os.Stat(scriptName)
	suite.True(os.IsNotExist(err), "File should be deleted")
}
