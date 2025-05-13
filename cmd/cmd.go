package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func NewCommand() (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "random-fit",
		Short: "Random Fit",
		Long:  `Random Fit generates random training sessions.`,
	}

	// definition represents the runEvent command
	var definition = &cobra.Command{
		Use:   "definition",
		Short: "Definition management",
		Long:  `Manage definitions: combination, plan.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	var code = &cobra.Command{
		Use:   "code",
		Short: "Code tools.",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	var codeGenerator = &cobra.Command{
		Use:   "generate",
		Short: "Generate code.",
		Run: func(cmd *cobra.Command, args []string) {
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
		// TODO Handle error
		os.Exit(1)
	}

	err = rootCmd.Execute()
	if err != nil {
		// TODO Handle error
		os.Exit(1)
	}

}
