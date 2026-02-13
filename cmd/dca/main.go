package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "dca",
	Short: "Discord CLI for Agentic Workflows",
	Long: `dca - A unified CLI tool for interacting with Discord,
designed for AI agents and automation workflows.

All commands output JSON for easy parsing by LLMs and scripts.`,
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/dca/config.json)")
	rootCmd.PersistentFlags().String("token", "", "Discord bot token (overrides config)")
	rootCmd.PersistentFlags().Bool("output-pretty", false, "Pretty print JSON output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
