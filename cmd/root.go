package cmd

import (
	"github.com/vcsfrl/random-fit/cmd/config"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger zerolog.Logger
var loggerOutput zerolog.ConsoleWriter
var appConfig *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "random-fit",
	Short: "Random Fit",
	Long:  `Random Fit generates random training sessions.`,
}

// apiCmd represents the runEvent command
var apiCmd = &cobra.Command{
	Use:   "combination",
	Short: "Combination management",
	Long:  `Manage combination related entities.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var exampleCmd = &cobra.Command{
	Use:   "coder",
	Short: "Code tools.",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	// Init logger.
	loggerOutput = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger = zerolog.New(loggerOutput).With().Timestamp().Logger()
	logger.Info().Msg("Logger initialised.")

	viper.SetConfigName("random-fit_config")
	viper.SetEnvPrefix("RF")

	if err := bindEnvConfig(apiCmd); err != nil {
		logger.Error().Err(err).Msg("Bind monitor config.")
		os.Exit(1)
	}

	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(exampleCmd)

	// Init config.
	appConfig = buildConfig(logger)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
