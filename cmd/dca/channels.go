package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulfschnabel/dca/internal/config"
	"github.com/ulfschnabel/dca/internal/discord"
	"github.com/ulfschnabel/dca/internal/output"
)

var channelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "Channel operations",
	Long:  "List channels and get message history",
}

var channelsListCmd = &cobra.Command{
	Use:   "list <server-id>",
	Short: "List channels in a server",
	Long:  "List all channels in a Discord server",
	Args:  cobra.ExactArgs(1),
	RunE:  runChannelsList,
}

var channelsHistoryCmd = &cobra.Command{
	Use:   "history <channel-id>",
	Short: "Get message history",
	Long:  "Get recent messages from a channel",
	Args:  cobra.ExactArgs(1),
	RunE:  runChannelsHistory,
}

func init() {
	rootCmd.AddCommand(channelsCmd)
	channelsCmd.AddCommand(channelsListCmd)
	channelsCmd.AddCommand(channelsHistoryCmd)

	channelsHistoryCmd.Flags().Int("limit", 10, "Number of messages to retrieve (max 100)")
}

func runChannelsList(cmd *cobra.Command, args []string) error {
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

	// List channels
	channels, err := client.ListChannels(serverID)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"channels": channels,
		"count":    len(channels),
	}, pretty)
}

func runChannelsHistory(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	channelID := args[0]
	limit, _ := cmd.Flags().GetInt("limit")

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

	// Get messages
	messages, err := client.GetMessages(channelID, limit)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"messages": messages,
		"count":    len(messages),
	}, pretty)
}
