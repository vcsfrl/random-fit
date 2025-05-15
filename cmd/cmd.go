package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vcsfrl/random-fit/cmd/internal"
	"log"
)

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
					if combinationDefinition, err := NewCombinationDefinition(cmd, args, NewConfig()); err == nil {
						combinationDefinition.List()
					}
				},
			}

			var newCombination = &cobra.Command{
				Use:   "new",
				Short: "New Combination Definition",
				Run: func(cmd *cobra.Command, args []string) {
					if combinationDefinition, err := NewCombinationDefinition(cmd, args, NewConfig()); err == nil {
						combinationDefinition.New()
					}
				},
			}

			var editCombination = &cobra.Command{
				Use:   "edit",
				Short: "Edit Combination Definition",
				Run: func(cmd *cobra.Command, args []string) {
					if combinationDefinition, err := NewCombinationDefinition(cmd, args, NewConfig()); err == nil {
						combinationDefinition.Edit()
					}
				},
			}

			var deleteCombination = &cobra.Command{
				Use:   "delete",
				Short: "Delete Combination Definition",
				Run: func(cmd *cobra.Command, args []string) {
					if combinationDefinition, err := NewCombinationDefinition(cmd, args, NewConfig()); err == nil {
						combinationDefinition.Delete()
					}
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

		// Plan Definition
		{
			var plan = &cobra.Command{
				Use:   "plan",
				Short: "Plan Definition management",
				Long:  `Manage plan definitions: list, new, edit, delete.`,
				Run: func(cmd *cobra.Command, args []string) {
					if planDefinition, err := NewPlanDefinition(cmd, args, NewConfig()); err == nil {
						planDefinition.List()
					}
				},
			}

			var newPlan = &cobra.Command{
				Use:   "new",
				Short: "New Plan Definition",
				Run: func(cmd *cobra.Command, args []string) {
					if planDefinition, err := NewPlanDefinition(cmd, args, NewConfig()); err == nil {
						planDefinition.New()
					}
				},
			}

			newPlan.Flags().String("name", "", "")

			plan.AddCommand(newPlan)
			definition.AddCommand(plan)
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
