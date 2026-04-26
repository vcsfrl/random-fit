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

const shutdownSleep = 200

func NewCommand() (*cobra.Command, error) {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "random-fit",
		Short: "Random Fit",
		Long:  `Random Fit generates random training sessions.`,
	}

	// Definition
	definitionCmd(rootCmd)
	generateCmd(rootCmd)
	codeCmd(rootCmd)
	runInteractiveCmd(rootCmd)

	viper.SetConfigName("random-fit_config")
	viper.SetEnvPrefix("RF")

	if err := BindEnvConfig(rootCmd); err != nil {
		return nil, err
	}

	return rootCmd, nil
}

func codeCmd(rootCmd *cobra.Command) {
	var code = &cobra.Command{
		Use:   "code",
		Short: "Code tools.",
	}

	var codeGenerator = &cobra.Command{
		Use:   "generate",
		Short: "Generate code.",
		Run: func(cmd *cobra.Command, _ []string) {
			service.GenerateCode(cmd, NewConfig().Config)
		},
	}

	code.AddCommand(codeGenerator)
	rootCmd.AddCommand(code)
}

func definitionCmd(rootCmd *cobra.Command) {
	var definition = &cobra.Command{
		Use:   "definition",
		Short: "Definition management",
		Long:  `Manage definitions: combination, plan.`,
		Run: func(_ *cobra.Command, _ []string) {
		},
	}

	// Combination Definition
	combinationDefinitionCmd(definition)
	planDefinitionCmd(definition)

	rootCmd.AddCommand(definition)
}

func generateCmd(rootCmd *cobra.Command) {
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
	generateCombination.Flags().Int("workers", defaultWorkers, "Number of export workers")
	generate.AddCommand(generateCombination)
	rootCmd.AddCommand(generate)
}

// DefinitionHandler defines the common interface for definition management operations.
type DefinitionHandler interface {
	New()
	Edit()
	Delete()
	List()
}

type definitionHandlerFactory func(*cobra.Command, []string, *Config) (DefinitionHandler, error)

func registerDefinitionSubcommands(
	parent *cobra.Command, use string, short string, long string, factory definitionHandlerFactory,
) {
	var root = &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			if handler, err := factory(cmd, args, NewConfig()); err == nil {
				handler.List()
			}
		},
	}

	var newCmd = &cobra.Command{
		Use:   "new",
		Short: "New " + short,
		Run: func(cmd *cobra.Command, args []string) {
			if handler, err := factory(cmd, args, NewConfig()); err == nil {
				handler.New()
			}
		},
	}

	var editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit " + short,
		Run: func(cmd *cobra.Command, args []string) {
			if handler, err := factory(cmd, args, NewConfig()); err == nil {
				handler.Edit()
			}
		},
	}

	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete " + short,
		Run: func(cmd *cobra.Command, args []string) {
			if handler, err := factory(cmd, args, NewConfig()); err == nil {
				handler.Delete()
			}
		},
	}

	newCmd.Flags().String("name", "", "")
	editCmd.Flags().String("name", "", "")
	deleteCmd.Flags().String("name", "", "")

	root.AddCommand(newCmd)
	root.AddCommand(editCmd)
	root.AddCommand(deleteCmd)
	parent.AddCommand(root)
}

//nolint:ireturn
func newCombinationHandler(cmd *cobra.Command, args []string, conf *Config) (DefinitionHandler, error) {
	return NewCombinationDefinition(cmd, args, conf)
}

//nolint:ireturn
func newPlanHandler(cmd *cobra.Command, args []string, conf *Config) (DefinitionHandler, error) {
	return NewPlanDefinition(cmd, args, conf)
}

func planDefinitionCmd(definition *cobra.Command) {
	registerDefinitionSubcommands(definition, "plan",
		"Plan Definition management",
		"Manage plan definitions: list, new, edit, delete.",
		newPlanHandler,
	)
}

func combinationDefinitionCmd(definition *cobra.Command) {
	registerDefinitionSubcommands(definition, "combination",
		"Combination Definition management",
		"Manage combination definitions: list, new, edit, delete.",
		newCombinationHandler,
	)
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

		time.Sleep(shutdownSleep * time.Millisecond)
		workGroup.Done()
		os.Exit(0)
	}()

	rootCmd, err := NewCommand()
	if err != nil {
		log.Printf("Error creating command: %v", err)

		return
	}

	err = rootCmd.ExecuteContext(ctx)
	if err != nil {
		log.Printf("Error executing command: %v", err)

		return
	}
}
