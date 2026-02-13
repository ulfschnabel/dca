package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ulfschnabel/dca/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management",
	Long:  "Manage dca configuration including bot token and settings",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long:  "Create a new configuration file with interactive prompts",
	RunE:  runConfigInit,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration (token will be masked)",
	RunE:  runConfigShow,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🤖 dca Configuration Setup")
	fmt.Println()
	fmt.Println("⚠️  This tool uses your Discord user token (personal automation)")
	fmt.Println()
	fmt.Println("📝 To get your Discord user token:")
	fmt.Println("   1. Open Discord in browser (web.discord.com)")
	fmt.Println("   2. Press F12 for Developer Tools")
	fmt.Println("   3. Go to Network tab")
	fmt.Println("   4. Send any message")
	fmt.Println("   5. Find 'authorization' header in requests")
	fmt.Println()
	fmt.Print("Discord User Token: ")
	userToken, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read user token: %w", err)
	}
	userToken = strings.TrimSpace(userToken)

	if userToken == "" {
		return fmt.Errorf("user token is required")
	}

	// Get approval setting
	fmt.Print("Require approval for write operations? [Y/n]: ")
	approvalInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read approval setting: %w", err)
	}
	approvalInput = strings.TrimSpace(strings.ToLower(approvalInput))
	requireApproval := approvalInput != "n" && approvalInput != "no"

	// Create config
	cfg := &config.Config{
		UserToken:       userToken,
		RequireApproval: requireApproval,
	}

	// Save config
	cfgPath := cfgFile
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath()
	}

	if err := config.Save(cfg, cfgPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("\n✅ Configuration saved to: %s\n", cfgPath)
	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfgPath := cfgFile
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath()
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}

	// Mask token
	maskedToken := cfg.UserToken
	if len(maskedToken) > 8 {
		maskedToken = maskedToken[:4] + "..." + maskedToken[len(maskedToken)-4:]
	}

	fmt.Printf("Config file: %s\n", cfgPath)
	fmt.Printf("User Token: %s\n", maskedToken)
	fmt.Printf("Require Approval: %v\n", cfg.RequireApproval)

	return nil
}
