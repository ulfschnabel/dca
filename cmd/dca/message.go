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

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Message operations",
	Long:  "Send, reply to, and manage messages",
}

var messageSendCmd = &cobra.Command{
	Use:   "send <channel-id> <message>",
	Short: "Send a message",
	Long:  "Send a message to a channel (requires approval unless --dry-run)",
	Args:  cobra.ExactArgs(2),
	RunE:  runMessageSend,
}

var messageReplyCmd = &cobra.Command{
	Use:   "reply <channel-id> <message-id> <message>",
	Short: "Reply to a message",
	Long:  "Reply to a specific message (requires approval unless --dry-run)",
	Args:  cobra.ExactArgs(3),
	RunE:  runMessageReply,
}

var messageEditCmd = &cobra.Command{
	Use:   "edit <channel-id> <message-id> <new-message>",
	Short: "Edit your message",
	Long:  "Edit one of your messages (requires approval unless --dry-run)",
	Args:  cobra.ExactArgs(3),
	RunE:  runMessageEdit,
}

var messageDeleteCmd = &cobra.Command{
	Use:   "delete <channel-id> <message-id>",
	Short: "Delete your message",
	Long:  "Delete one of your messages (requires approval unless --dry-run)",
	Args:  cobra.ExactArgs(2),
	RunE:  runMessageDelete,
}

func init() {
	rootCmd.AddCommand(messageCmd)
	messageCmd.AddCommand(messageSendCmd)
	messageCmd.AddCommand(messageReplyCmd)
	messageCmd.AddCommand(messageEditCmd)
	messageCmd.AddCommand(messageDeleteCmd)

	messageSendCmd.Flags().Bool("dry-run", false, "Show what would be sent without actually sending")
	messageReplyCmd.Flags().Bool("dry-run", false, "Show what would be sent without actually sending")
	messageEditCmd.Flags().Bool("dry-run", false, "Show what would be changed without actually changing")
	messageDeleteCmd.Flags().Bool("dry-run", false, "Show what would be deleted without actually deleting")
}

func runMessageSend(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	channelID := args[0]
	content := args[1]

	// Load config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	// Dry run - just show what would be sent
	if dryRun {
		return output.PrintSuccess(map[string]interface{}{
			"action":     "send_message",
			"channel_id": channelID,
			"content":    content,
			"dry_run":    true,
		}, pretty)
	}

	// Check approval requirement
	if cfg.RequireApproval {
		fmt.Printf("📝 Send message to channel %s:\n", channelID)
		fmt.Printf("   \"%s\"\n\n", content)
		fmt.Print("Proceed? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return output.PrintError(fmt.Errorf("failed to read response: %w", err), pretty)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			return output.PrintSuccess(map[string]interface{}{
				"action":    "send_message",
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

	// Send message
	msg, err := client.SendMessage(channelID, content)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(msg, pretty)
}

func runMessageReply(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	channelID := args[0]
	messageID := args[1]
	content := args[2]

	// Load config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	// Dry run
	if dryRun {
		return output.PrintSuccess(map[string]interface{}{
			"action":     "reply_message",
			"channel_id": channelID,
			"message_id": messageID,
			"content":    content,
			"dry_run":    true,
		}, pretty)
	}

	// Check approval requirement
	if cfg.RequireApproval {
		fmt.Printf("📝 Reply to message %s in channel %s:\n", messageID, channelID)
		fmt.Printf("   \"%s\"\n\n", content)
		fmt.Print("Proceed? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return output.PrintError(fmt.Errorf("failed to read response: %w", err), pretty)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			return output.PrintSuccess(map[string]interface{}{
				"action":    "reply_message",
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

	// Reply to message
	msg, err := client.ReplyToMessage(channelID, messageID, content)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(msg, pretty)
}

func runMessageEdit(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	channelID := args[0]
	messageID := args[1]
	newContent := args[2]

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

	// Get original message for approval prompt
	originalMsg, err := client.GetMessage(channelID, messageID)
	if err != nil {
		return output.PrintError(fmt.Errorf("failed to get original message: %w", err), pretty)
	}

	// Dry run
	if dryRun {
		return output.PrintSuccess(map[string]interface{}{
			"action":           "edit_message",
			"channel_id":       channelID,
			"message_id":       messageID,
			"original_content": originalMsg.Content,
			"new_content":      newContent,
			"dry_run":          true,
		}, pretty)
	}

	// Check approval requirement
	if cfg.RequireApproval {
		fmt.Printf("📝 Edit message %s in channel %s:\n", messageID, channelID)
		fmt.Printf("   Old: \"%s\"\n", originalMsg.Content)
		fmt.Printf("   New: \"%s\"\n\n", newContent)
		fmt.Print("Proceed? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return output.PrintError(fmt.Errorf("failed to read response: %w", err), pretty)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			return output.PrintSuccess(map[string]interface{}{
				"action":    "edit_message",
				"cancelled": true,
			}, pretty)
		}
	}

	// Edit message
	msg, err := client.EditMessage(channelID, messageID, newContent)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(msg, pretty)
}

func runMessageDelete(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	channelID := args[0]
	messageID := args[1]

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

	// Get original message for approval prompt
	originalMsg, err := client.GetMessage(channelID, messageID)
	if err != nil {
		return output.PrintError(fmt.Errorf("failed to get message: %w", err), pretty)
	}

	// Dry run
	if dryRun {
		return output.PrintSuccess(map[string]interface{}{
			"action":     "delete_message",
			"channel_id": channelID,
			"message_id": messageID,
			"content":    originalMsg.Content,
			"dry_run":    true,
		}, pretty)
	}

	// Check approval requirement
	if cfg.RequireApproval {
		fmt.Printf("🗑️  Delete message %s in channel %s:\n", messageID, channelID)
		fmt.Printf("   \"%s\"\n\n", originalMsg.Content)
		fmt.Print("Proceed? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return output.PrintError(fmt.Errorf("failed to read response: %w", err), pretty)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			return output.PrintSuccess(map[string]interface{}{
				"action":    "delete_message",
				"cancelled": true,
			}, pretty)
		}
	}

	// Delete message
	err = client.DeleteMessage(channelID, messageID)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"action":     "delete_message",
		"channel_id": channelID,
		"message_id": messageID,
		"deleted":    true,
	}, pretty)
}
