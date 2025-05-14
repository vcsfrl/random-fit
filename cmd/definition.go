package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/cmd/internal"
	"os"
	"os/exec"
)

type BaseDefinition struct {
	cmd  *cobra.Command
	args []string
	conf *internal.Config
}

func (b *BaseDefinition) editScript(scriptName string, filetype string) error {
	if os.Getenv("EDITOR") == "" {
		return errNoEnvEditor
	}
	cmd := exec.Command(os.Getenv("EDITOR"), "-filetype", filetype, scriptName)
	cmd.Stdin = b.cmd.InOrStdin()
	cmd.Stdout = b.cmd.OutOrStdout()
	cmd.Stderr = b.cmd.ErrOrStderr()

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (b *BaseDefinition) getNameArg() string {
	name := ""
	if len(b.args) > 0 {
		name = b.args[0]
	}

	if name == "" {
		name, _ = b.cmd.Flags().GetString("name")
	}

	return name
}

func (b *BaseDefinition) createFolder(folder string) error {
	return createFolder(folder)
}

func createFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
}
