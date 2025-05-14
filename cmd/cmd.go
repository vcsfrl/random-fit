package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vcsfrl/random-fit/cmd/internal"
	"log"
)

var errNoEnvEditor = fmt.Errorf("EDITOR environment variable is not set")

var msgNameMissing = "Name is required."
var msgCombinationDefinition = "Combination Definition"
var msqCreate = "Create"
var msqCreated = "Created"
var msqEditScript = "Editing script"

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
			NewCombinationDefinition(cmd, args, conf).New()
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
