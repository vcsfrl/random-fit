package app

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	TracePort       string
	DebuggerPort    string
	DataFolder      string
	BaseFolder      string
	K8sSharedFolder string
	Editor          string
}

func NewConfig(logger zerolog.Logger) *Config {
	logger.Info().Msg("Initialize api app.")

	var newConfig Config
	newConfig.TracePort = viper.Get("tracePort").(string)
	newConfig.DebuggerPort = viper.Get("debuggerPort").(string)
	newConfig.DataFolder = viper.Get("dataFolder").(string)
	newConfig.BaseFolder = viper.Get("baseFolder").(string)
	newConfig.K8sSharedFolder = viper.Get("k8sSharedFolder").(string)
	newConfig.Editor = viper.Get("editor").(string)

	return &newConfig
}

func BindEnvConfig(command *cobra.Command) error {
	command.Flags().String("trace-port", "40021", "Trace port")
	if err := viper.BindPFlag("tracePort", command.Flags().Lookup("trace-port")); err != nil {
		return err
	}
	if err := viper.BindEnv("tracePort", "RF_TRACE_PORT"); err != nil {
		return err
	}

	command.Flags().String("debugger-port", "40022", "Debugger port")
	if err := viper.BindPFlag("debuggerPort", command.Flags().Lookup("debugger-port")); err != nil {
		return err
	}
	if err := viper.BindEnv("debuggerPort", "RF_DEBUGGER_PORT"); err != nil {
		return err
	}

	command.Flags().String("data-folder", "/srv/random-fit/data", "Data folder")
	if err := viper.BindPFlag("dataFolder", command.Flags().Lookup("data-folder")); err != nil {
		return err
	}
	if err := viper.BindEnv("dataFolder", "RF_DATA_FOLDER"); err != nil {
		return err
	}

	command.Flags().String("base-folder", "/srv/random-fit", "Base folder")
	if err := viper.BindPFlag("baseFolder", command.Flags().Lookup("base-folder")); err != nil {
		return err
	}
	if err := viper.BindEnv("baseFolder", "RF_BASE_FOLDER"); err != nil {
		return err
	}

	command.Flags().String("k8s-shared-folder", "k8s-shared", "K8s shared folder")
	if err := viper.BindPFlag("k8sSharedFolder", command.Flags().Lookup("k8s-shared-folder")); err != nil {
		return err
	}
	if err := viper.BindEnv("k8sSharedFolder", "RF_K8S_SHARED_FOLDER"); err != nil {
		return err
	}

	command.Flags().String("editor", "micro", "Editor")
	if err := viper.BindPFlag("editor", command.Flags().Lookup("editor")); err != nil {
		return err
	}
	if err := viper.BindEnv("editor", "EDITOR"); err != nil {
		return err
	}

	return nil
}
