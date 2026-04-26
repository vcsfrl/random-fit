package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/vcsfrl/random-fit/internal/platform/fs"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	// Initialize translations.
	_ "github.com/vcsfrl/random-fit/cmd/translations"
)

var ErrNoEnvEditor = errors.New("EDITOR environment variable is not set")
var ErrInvalidDefinitionName = errors.New(
	"invalid definition name: must contain only letters, digits, hyphens, and underscores",
)

type BaseHandler struct {
	cmd  *cobra.Command
	args []string
	conf *Config

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

var validNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func (b *BaseHandler) validateName(name string) error {
	if !validNamePattern.MatchString(name) {
		return ErrInvalidDefinitionName
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
	folders := []struct {
		path string
		name string
	}{
		{b.conf.DefinitionFolder(), "definition"},
		{b.conf.PlanFolder(), "plan"},
		{b.conf.CombinationFolder(), "combination"},
		{b.conf.StorageFolder(), "storage"},
	}

	for _, folder := range folders {
		if err := b.createFolder(folder.path); err != nil {
			b.cmd.PrintErrln(b.printer.Sprintf("Error creating %s folder:", folder.name), err)

			return err
		}
	}

	return nil
}
