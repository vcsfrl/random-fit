package shell

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestShell(t *testing.T) {
	suite.Run(t, new(ShellSuite))
}

type ShellSuite struct {
	suite.Suite

	shell       *Shell
	input       *Buffer
	inputWriter *Buffer
	output      *Buffer
	errors      *Buffer
}

func (suite *ShellSuite) SetupTest() {
	suite.input = &Buffer{}
	suite.output = &Buffer{}
	suite.errors = &Buffer{}
	suite.inputWriter = &Buffer{}

	// Create a new shell instance
	suite.shell = BuildNew()

	// Check if the shell instance is not nil
	suite.NotNil(suite.shell)

	suite.NotNil(suite.shell.stdin)
	suite.NotNil(suite.shell.stdout)
	suite.NotNil(suite.shell.stderr)

	// Create a new shell instance with custom input and output
	suite.shell.stdin = suite.input
	suite.shell.stdinWriter = suite.inputWriter
	suite.shell.stdout = suite.output
	suite.shell.stderr = suite.errors
	suite.shell.Init()

}

func (suite *ShellSuite) TearDownTest() {
	// Close the shell instance
	if suite.shell != nil {
		_ = suite.shell.Close()
	}
}

func (suite *ShellSuite) TestNew() {

	// Check if the shell instance is not nil
	suite.NotNil(suite.shell)

	// Check if the shell instance has the expected properties
	suite.NotNil(suite.shell.definitionFolder)
	suite.NotNil(suite.shell.planFolder)
	suite.NotNil(suite.shell.storageFolder)
	suite.NotNil(suite.shell.combinationFolder)
	suite.NotNil(suite.shell.ctxCancel)
	suite.NotNil(suite.shell.ctx)
}

func (suite *ShellSuite) TestRun() {
	// Run the shell instance
	suite.shell.Run()

	output := suite.output.String()

	suite.Contains(output, welcomeMessage)
	suite.Contains(output, separator)
}

func (suite *ShellSuite) TestHelp() {
	// Run the shell instance
	suite.shell.Run()

	// Check if the help command is available
	suite.shell.RunCommand("help")

	output := suite.output.String()

	suite.Contains(output, "Commands:")
}

type Buffer struct {
	buf *bytes.Buffer
}

func (cb *Buffer) Write(p []byte) (n int, err error) {
	if cb.buf == nil {
		cb.buf = &bytes.Buffer{}
	}
	return cb.buf.Write(p)
}

func (cb *Buffer) Read(p []byte) (n int, err error) {
	if cb.buf == nil {
		return 0, nil
	}
	return cb.buf.Read(p)
}

func (cb *Buffer) Close() error {
	if cb.buf != nil {
		cb.buf.Reset()
	}
	return nil
}

func (cb *Buffer) String() string {
	if cb.buf != nil {
		return cb.buf.String()
	}
	return ""
}
