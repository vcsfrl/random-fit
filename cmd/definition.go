package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/platform/fs"
	"github.com/vcsfrl/random-fit/internal/service"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os/exec"

	// Initialize translations.
	_ "github.com/vcsfrl/random-fit/cmd/translations"
)

var ErrNoEnvEditor = errors.New("EDITOR environment variable is not set")

type BaseHandler struct {
	cmd  *cobra.Command
	args []string
	conf *service.Config

	printer *message.Printer
}

func (b *BaseHandler) editScript(scriptName string, filetype string) error {
	if b.conf.Editor == "" || b.conf.Editor == "-" {
		return ErrNoEnvEditor
	}

	cmd := exec.Command(b.conf.Editor, "-filetype", filetype, scriptName) //nolint:gosec
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
	err := fs.CreateFolder(folder)
	if err != nil {
		return fmt.Errorf("error creating folder %s: %w", folder, err)
	}

	return nil
}

func (b *BaseHandler) initTranslations() {
	var lang language.Tag

	switch b.conf.Locale {
	case "en_US.UTF-8":
		lang = language.MustParse("en-US")
	default:
		lang = language.MustParse("en-US")
	}

	b.printer = message.NewPrinter(lang)
}

func (b *BaseHandler) initFolders() error {
	err := b.createFolder(b.conf.DefinitionFolder())
	if err != nil {
		b.cmd.PrintErrln(b.printer.Sprintf("Error creating definition folder:"), err)

		return err
	}

	err = b.createFolder(b.conf.PlanFolder())
	if err != nil {
		b.cmd.PrintErrln(b.printer.Sprintf("Error creating plan folder:"), err)

		return err
	}

	err = b.createFolder(b.conf.CombinationFolder())
	if err != nil {
		b.cmd.PrintErrln(b.printer.Sprintf("Error creating combination folder:"), err)

		return err
	}

	err = b.createFolder(b.conf.StorageFolder())
	if err != nil {
		b.cmd.PrintErrln(b.printer.Sprintf("Error creating storage folder:"), err)

		return err
	}

	return nil
}
