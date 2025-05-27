package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vcsfrl/random-fit/internal/service"
)

func NewConfig() *service.Config {
	var newConfig service.Config
	newConfig.TracePort = viper.GetString("tracePort")
	newConfig.DebuggerPort = viper.GetString("debuggerPort")
	newConfig.DebugChartPort = viper.GetString("debugChartPort")
	newConfig.DataFolder = viper.GetString("dataFolder")
	newConfig.BaseFolder = viper.GetString("baseFolder")
	newConfig.K8sSharedFolder = viper.GetString("k8sSharedFolder")
	newConfig.Editor = viper.GetString("editor")
	newConfig.Locale = viper.GetString("locale")

	return &newConfig
}

func BindEnvConfig(command *cobra.Command) error {
	command.Flags().String("trace-port", "40021", "Trace port")

	if err := viper.BindPFlag("tracePort", command.Flags().Lookup("trace-port")); err != nil {
		return fmt.Errorf("error binding flag: %w", err)
	}

	if err := viper.BindEnv("tracePort", "RF_TRACE_PORT"); err != nil {
		return fmt.Errorf("error binding env: %w", err)
	}

	command.Flags().String("debugger-port", "40022", "Debugger port")

	if err := viper.BindPFlag("debuggerPort", command.Flags().Lookup("debugger-port")); err != nil {
		return fmt.Errorf("error binding flag: %w", err)
	}

	if err := viper.BindEnv("debuggerPort", "RF_DEBUGGER_PORT"); err != nil {
		return fmt.Errorf("error binding env: %w", err)
	}

	command.Flags().String("debug-chart-port", "40023", "Debug chart port")

	if err := viper.BindPFlag("debugChartPort", command.Flags().Lookup("debug-chart-port")); err != nil {
		return fmt.Errorf("error binding flag: %w", err)
	}

	if err := viper.BindEnv("debugChartPort", "RF_DEBUG_CHART_PORT"); err != nil {
		return fmt.Errorf("error binding env: %w", err)
	}

	command.Flags().String("data-folder", "/srv/random-fit/data", "Data folder")

	if err := viper.BindPFlag("dataFolder", command.Flags().Lookup("data-folder")); err != nil {
		return fmt.Errorf("error binding flag: %w", err)
	}

	if err := viper.BindEnv("dataFolder", "RF_DATA_FOLDER"); err != nil {
		return fmt.Errorf("error binding env: %w", err)
	}

	command.Flags().String("base-folder", "/srv/random-fit", "Base folder")

	if err := viper.BindPFlag("baseFolder", command.Flags().Lookup("base-folder")); err != nil {
		return fmt.Errorf("error binding flag: %w", err)
	}

	if err := viper.BindEnv("baseFolder", "RF_BASE_FOLDER"); err != nil {
		return fmt.Errorf("error binding env: %w", err)
	}

	command.Flags().String("k8s-shared-folder", "k8s-shared", "K8s shared folder")

	if err := viper.BindPFlag("k8sSharedFolder", command.Flags().Lookup("k8s-shared-folder")); err != nil {
		return fmt.Errorf("error binding flag: %w", err)
	}

	if err := viper.BindEnv("k8sSharedFolder", "RF_K8S_SHARED_FOLDER"); err != nil {
		return fmt.Errorf("error binding env: %w", err)
	}

	command.Flags().String("editor", "micro", "Editor")

	if err := viper.BindPFlag("editor", command.Flags().Lookup("editor")); err != nil {
		return fmt.Errorf("error binding flag: %w", err)
	}

	if err := viper.BindEnv("editor", "EDITOR"); err != nil {
		return fmt.Errorf("error binding env: %w", err)
	}

	command.Flags().String("locale", "en-US", "Locale")

	if err := viper.BindPFlag("locale", command.Flags().Lookup("locale")); err != nil {
		return fmt.Errorf("error binding flag: %w", err)
	}

	if err := viper.BindEnv("locale", "RF_LOCALE"); err != nil {
		return fmt.Errorf("error binding env: %w", err)
	}

	return nil
}
