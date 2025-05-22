package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/service"
	"os"
	"os/exec"
)

type BaseHandler struct {
	cmd  *cobra.Command
	args []string
	conf *service.Config
}

func (b *BaseHandler) editScript(scriptName string, filetype string) error {
	if os.Getenv("EDITOR") == "" {
		return ErrNoEnvEditor
	}
	cmd := exec.Command(os.Getenv("EDITOR"), "-filetype", filetype, scriptName)
	cmd.Stdin = b.cmd.InOrStdin()
	cmd.Stdout = b.cmd.OutOrStdout()
	cmd.Stderr = b.cmd.ErrOrStderr()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting editor: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for editor: %w", err)
	}

	return nil
}

func (b *BaseHandler) getArg(position int, flagName string) string {
	arg := ""

	if len(b.args) > 0 && len(b.args)-1 >= position {
		arg = b.args[position]
	}

	if arg == "" {
		arg, _ = b.cmd.Flags().GetString(flagName)
	}

	return arg
}

func (b *BaseHandler) createFolder(folder string) error {
	return CreateFolder(folder)
}

func (b *BaseHandler) initFolders() error {
	err := b.createFolder(b.conf.DefinitionFolder())
	if err != nil {
		b.cmd.PrintErrln("Error creating definition folder: ", err)
		return err
	}

	err = b.createFolder(b.conf.PlanFolder())
	if err != nil {
		b.cmd.PrintErrln("Error creating plan folder: ", err)
		return err
	}

	err = b.createFolder(b.conf.CombinationFolder())
	if err != nil {
		b.cmd.PrintErrln("Error creating combination folder: ", err)
		return err
	}

	err = b.createFolder(b.conf.StorageFolder())
	if err != nil {
		b.cmd.PrintErrln("Error creating storage folder: ", err)
		return err
	}

	return nil
}

func CreateFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return fmt.Errorf("error creating folder %s: %w", folder, err)
		}
	}

	return nil
}
