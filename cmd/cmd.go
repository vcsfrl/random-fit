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
var msgList = "List"
var msgCreate = "Create"
var msgEdit = "Edit"
var msgDelete = "Delete"
var msgDone = "DONE:"
var msgEditScript = "Editing script"
var msgRemoveScript = "Removing script"
var msgNoItemsFound = "No items found!"

func NewCommand() (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "random-fit",
		Short: "Random Fit",
		Long:  `Random Fit generates random training sessions.`,
	}

	// Definition
	{
		var definition = &cobra.Command{
			Use:   "definition",
			Short: "Definition management",
			Long:  `Manage definitions: combination, plan.`,
			Run: func(cmd *cobra.Command, args []string) {

			},
		}

		// Combination Definition
		{
			var combination = &cobra.Command{
				Use:   "combination",
				Short: "Combination Definition management",
				Long:  `Manage combination definitions: list, new, edit, delete.`,
				Run: func(cmd *cobra.Command, args []string) {
					conf := NewConfig()
					NewCombinationDefinition(cmd, args, conf).List()
				},
			}

			var newCombination = &cobra.Command{
				Use:   "new",
				Short: "New Combination Definition",
				Run: func(cmd *cobra.Command, args []string) {
					conf := NewConfig()
					NewCombinationDefinition(cmd, args, conf).New()
				},
			}

			var editCombination = &cobra.Command{
				Use:   "edit",
				Short: "Edit Combination Definition",
				Run: func(cmd *cobra.Command, args []string) {
					conf := NewConfig()
					NewCombinationDefinition(cmd, args, conf).Edit()
				},
			}

			var deleteCombination = &cobra.Command{
				Use:   "delete",
				Short: "Delete Combination Definition",
				Run: func(cmd *cobra.Command, args []string) {
					conf := NewConfig()
					NewCombinationDefinition(cmd, args, conf).Delete()
				},
			}

			newCombination.Flags().String("name", "", "")
			editCombination.Flags().String("name", "", "")
			deleteCombination.Flags().String("name", "", "")

			combination.AddCommand(newCombination)
			combination.AddCommand(editCombination)
			combination.AddCommand(deleteCombination)
			definition.AddCommand(combination)
		}

		rootCmd.AddCommand(definition)
	}

	{
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
		rootCmd.AddCommand(code)
	}

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
