package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulfschnabel/dca/internal/config"
	"github.com/ulfschnabel/dca/internal/discord"
	"github.com/ulfschnabel/dca/internal/output"
)

var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "Server operations",
	Long:  "List and get information about Discord servers (guilds)",
}

var serversListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all servers",
	Long:  "List all Discord servers you have access to",
	RunE:  runServersList,
}

var serversInfoCmd = &cobra.Command{
	Use:   "info <server-id>",
	Short: "Get server information",
	Long:  "Get detailed information about a specific server",
	Args:  cobra.ExactArgs(1),
	RunE:  runServersInfo,
}

func init() {
	rootCmd.AddCommand(serversCmd)
	serversCmd.AddCommand(serversListCmd)
	serversCmd.AddCommand(serversInfoCmd)
}

func runServersList(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")

	// Load config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	// Get token from flag or config
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		token = cfg.UserToken
	}

	if token == "" {
		return output.PrintError(fmt.Errorf("no token configured"), pretty)
	}

	// Create Discord client
	client, err := discord.New(token)
	if err != nil {
		return output.PrintError(err, pretty)
	}
	defer client.Close()

	// List guilds
	guilds, err := client.ListGuilds()
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"servers": guilds,
		"count":   len(guilds),
	}, pretty)
}

func runServersInfo(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	serverID := args[0]

	// Load config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	// Get token
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		token = cfg.UserToken
	}

	if token == "" {
		return output.PrintError(fmt.Errorf("no token configured"), pretty)
	}

	// Create Discord client
	client, err := discord.New(token)
	if err != nil {
		return output.PrintError(err, pretty)
	}
	defer client.Close()

	// Get guild info
	guild, err := client.GetGuild(serverID)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(guild, pretty)
}
