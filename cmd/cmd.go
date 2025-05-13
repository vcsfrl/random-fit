package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vcsfrl/random-fit/cmd/internal"
	"log"
	"os"
	"os/exec"
)

func NewCommand() (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "random-fit",
		Short: "Random Fit",
		Long:  `Random Fit generates random training sessions.`,
	}

	var definition = &cobra.Command{
		Use:   "definition",
		Short: "Definition management",
		Long:  `Manage definitions: combination, plan.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	var combination = &cobra.Command{
		Use:   "combination",
		Short: "Combination definition management",
	}

	var newCombination = &cobra.Command{
		Use:   "new",
		Short: "New Combination definition",
		Run: func(cmd *cobra.Command, args []string) {
			conf := NewConfig()

			name := ""
			if len(args) > 0 {
				name = args[0]
			}

			if name == "" {
				name, _ = cmd.Flags().GetString("name")
			}

			if name == "" {
				cmd.PrintErrln("Error: name is required.")
				return
			}

			cmd.Printf("Create new Combination definition: %s\n", name)
			if err := createFolder(conf.DefinitionFolder()); err != nil {
				cmd.PrintErrln("Error creating definition folder: ", err)
				return
			}

			definitionManager := internal.NewCombinationStarDefinitionManager(conf.DefinitionFolder())
			err := definitionManager.New(name)
			if err != nil {
				cmd.PrintErrln("Error: ", err)
				return
			}

			cmd.Printf("Combination definition '%s' created.\n", name)

			scriptName, err := definitionManager.GetScript(name)
			if err != nil {
				cmd.PrintErrln("Error getting script: ", err)
				return
			}

			if err := editScript(scriptName, "python", cmd); err != nil {
				cmd.PrintErrln("Error editing script: ", err)
				return
			}

		},
	}

	newCombination.Flags().String("name", "", "")

	definition.AddCommand(combination)
	combination.AddCommand(newCombination)

	var code = &cobra.Command{
		Use:   "code",
		Short: "Code tools.",
	}
	var codeGenerator = &cobra.Command{
		Use:   "generate",
		Short: "Generate code.",
		Run: func(cmd *cobra.Command, args []string) {
			internal.GenerateCode(cmd, NewConfig())
		},
	}

	code.AddCommand(codeGenerator)
	rootCmd.AddCommand(definition)
	rootCmd.AddCommand(code)

	viper.SetConfigName("random-fit_config")
	viper.SetEnvPrefix("RF")
	if err := BindEnvConfig(rootCmd); err != nil {
		return nil, err
	}

	return rootCmd, nil
}

func createFolder(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}

	return nil
}

func editScript(scriptName string, filetype string, cliCmd *cobra.Command) error {
	if os.Getenv("EDITOR") == "" {
		return fmt.Errorf("EDITOR environment variable is not set")
	}
	cmd := exec.Command(os.Getenv("EDITOR"), "-filetype", filetype, scriptName)
	cmd.Stdin = cliCmd.InOrStdin()
	cmd.Stdout = cliCmd.OutOrStdout()
	cmd.Stderr = cliCmd.ErrOrStderr()

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd, err := NewCommand()
	if err != nil {
		log.Fatal(err)
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

}
