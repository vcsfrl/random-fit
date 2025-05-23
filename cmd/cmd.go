package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vcsfrl/random-fit/internal/service"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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
			Run: func(_ *cobra.Command, _ []string) {
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

			var editPlan = &cobra.Command{
				Use:   "edit",
				Short: "Edit Plan Definition",
				Run: func(cmd *cobra.Command, args []string) {
					if planDefinition, err := NewPlanDefinition(cmd, args, NewConfig()); err == nil {
						planDefinition.Edit()
					}
				},
			}

			var deletePlan = &cobra.Command{
				Use:   "delete",
				Short: "Delete Plan Definition",
				Run: func(cmd *cobra.Command, args []string) {
					if planDefinition, err := NewPlanDefinition(cmd, args, NewConfig()); err == nil {
						planDefinition.Delete()
					}
				},
			}

			newPlan.Flags().String("name", "", "")
			editPlan.Flags().String("name", "", "")
			deletePlan.Flags().String("name", "", "")

			plan.AddCommand(newPlan)
			plan.AddCommand(editPlan)
			plan.AddCommand(deletePlan)
			definition.AddCommand(plan)
		}

		rootCmd.AddCommand(definition)
	}

	// Generate
	{
		var generate = &cobra.Command{
			Use:   "generate",
			Short: "Generate entities: combinations",
		}

		var generateCombination = &cobra.Command{
			Use:   "combination",
			Short: "Generate a combination",
			Long:  `Generate a combination from a combination definition and a plan definition.`,
			Run: func(cmd *cobra.Command, args []string) {
				if generator, err := NewGenerator(cmd, args, NewConfig()); err == nil {
					generator.Combination()
				}
			},
		}

		generateCombination.Flags().String("combination", "", "Combination definition name")
		generateCombination.Flags().String("plan", "", "Plan definition name")
		generate.AddCommand(generateCombination)
		rootCmd.AddCommand(generate)
	}

	{
		var code = &cobra.Command{
			Use:   "code",
			Short: "Code tools.",
		}
		var codeGenerator = &cobra.Command{
			Use:   "generate",
			Short: "Generate code.",
			Run: func(cmd *cobra.Command, _ []string) {
				service.GenerateCode(cmd, NewConfig())
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	workGroup := &sync.WaitGroup{}

	defer func() {
		stop()
		workGroup.Wait()
	}()

	workGroup.Add(1)

	go func() {
		<-ctx.Done()

		time.Sleep(200 * time.Millisecond)
		workGroup.Done()
		os.Exit(0)
	}()

	rootCmd, err := NewCommand()
	if err != nil {
		log.Fatal(err)
	}

	err = rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
