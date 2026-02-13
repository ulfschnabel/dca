package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ulfschnabel/dca/internal/config"
	"github.com/ulfschnabel/dca/internal/discord"
	"github.com/ulfschnabel/dca/internal/output"
)

var reactionCmd = &cobra.Command{
	Use:   "reaction",
	Short: "Reaction operations",
	Long:  "Add, remove, and list message reactions",
}

var reactionAddCmd = &cobra.Command{
	Use:   "add <channel-id> <message-id> <emoji>",
	Short: "Add a reaction",
	Long:  "Add an emoji reaction to a message (requires approval unless --dry-run)",
	Args:  cobra.ExactArgs(3),
	RunE:  runReactionAdd,
}

var reactionRemoveCmd = &cobra.Command{
	Use:   "remove <channel-id> <message-id> <emoji>",
	Short: "Remove a reaction",
	Long:  "Remove your emoji reaction from a message (requires approval unless --dry-run)",
	Args:  cobra.ExactArgs(3),
	RunE:  runReactionRemove,
}

func init() {
	rootCmd.AddCommand(reactionCmd)
	reactionCmd.AddCommand(reactionAddCmd)
	reactionCmd.AddCommand(reactionRemoveCmd)

	reactionAddCmd.Flags().Bool("dry-run", false, "Show what would be added without actually adding")
	reactionRemoveCmd.Flags().Bool("dry-run", false, "Show what would be removed without actually removing")
}

func runReactionAdd(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	channelID := args[0]
	messageID := args[1]
	emoji := args[2]

	// Load config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	// Dry run
	if dryRun {
		return output.PrintSuccess(map[string]interface{}{
			"action":     "add_reaction",
			"channel_id": channelID,
			"message_id": messageID,
			"emoji":      emoji,
			"dry_run":    true,
		}, pretty)
	}

	// Check approval requirement
	if cfg.RequireApproval {
		fmt.Printf("👍 Add reaction %s to message %s in channel %s\n\n", emoji, messageID, channelID)
		fmt.Print("Proceed? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return output.PrintError(fmt.Errorf("failed to read response: %w", err), pretty)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			return output.PrintSuccess(map[string]interface{}{
				"action":    "add_reaction",
				"cancelled": true,
			}, pretty)
		}
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

	// Add reaction
	err = client.AddReaction(channelID, messageID, emoji)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"action":     "add_reaction",
		"channel_id": channelID,
		"message_id": messageID,
		"emoji":      emoji,
		"added":      true,
	}, pretty)
}

func runReactionRemove(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	channelID := args[0]
	messageID := args[1]
	emoji := args[2]

	// Load config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	// Dry run
	if dryRun {
		return output.PrintSuccess(map[string]interface{}{
			"action":     "remove_reaction",
			"channel_id": channelID,
			"message_id": messageID,
			"emoji":      emoji,
			"dry_run":    true,
		}, pretty)
	}

	// Check approval requirement
	if cfg.RequireApproval {
		fmt.Printf("❌ Remove reaction %s from message %s in channel %s\n\n", emoji, messageID, channelID)
		fmt.Print("Proceed? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return output.PrintError(fmt.Errorf("failed to read response: %w", err), pretty)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			return output.PrintSuccess(map[string]interface{}{
				"action":    "remove_reaction",
				"cancelled": true,
			}, pretty)
		}
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

	// Remove reaction
	err = client.RemoveReaction(channelID, messageID, emoji)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"action":     "remove_reaction",
		"channel_id": channelID,
		"message_id": messageID,
		"emoji":      emoji,
		"removed":    true,
	}, pretty)
}
