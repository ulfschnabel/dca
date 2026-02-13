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

var dmCmd = &cobra.Command{
	Use:   "dm",
	Short: "Direct message operations",
	Long:  "Send direct messages to users",
}

var dmSendCmd = &cobra.Command{
	Use:   "send <username-or-id> <message>",
	Short: "Send a DM to a user",
	Long:  "Send a direct message to a user by username or ID (requires approval unless --dry-run)",
	Args:  cobra.ExactArgs(2),
	RunE:  runDMSend,
}

var dmHistoryCmd = &cobra.Command{
	Use:   "history <username-or-id>",
	Short: "Get DM history with a user",
	Long:  "Get recent direct messages with a user",
	Args:  cobra.ExactArgs(1),
	RunE:  runDMHistory,
}

var dmListCmd = &cobra.Command{
	Use:   "list",
	Short: "List DM conversations",
	Long:  "List all DM conversations sorted by recent activity",
	RunE:  runDMList,
}

func init() {
	rootCmd.AddCommand(dmCmd)
	dmCmd.AddCommand(dmSendCmd)
	dmCmd.AddCommand(dmHistoryCmd)
	dmCmd.AddCommand(dmListCmd)

	dmSendCmd.Flags().Bool("dry-run", false, "Show what would be sent without actually sending")
	dmHistoryCmd.Flags().Int("limit", 10, "Number of messages to retrieve")
	dmListCmd.Flags().Int("limit", 20, "Number of DM channels to show")
	dmListCmd.Flags().Bool("active-only", true, "Only show DMs with recent messages")
}

func runDMSend(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	userIdentifier := args[0]
	content := args[1]

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

	// Resolve username to user ID if needed
	userID := userIdentifier
	var username string

	// If input doesn't look like a numeric ID, treat it as a username
	if !isNumeric(userIdentifier) {
		fmt.Printf("🔍 Looking up user '%s'...\n", userIdentifier)
		user, err := client.FindUserByUsername(userIdentifier)
		if err != nil {
			return output.PrintError(err, pretty)
		}
		userID = user.ID
		username = user.Username
		fmt.Printf("✓ Found: %s (ID: %s)\n\n", username, userID)
	} else {
		username = userID // Use ID as display name if given directly
	}

	// Dry run - just show what would be sent
	if dryRun {
		return output.PrintSuccess(map[string]interface{}{
			"action":   "send_dm",
			"user_id":  userID,
			"username": username,
			"content":  content,
			"dry_run":  true,
		}, pretty)
	}

	// Check approval requirement
	if cfg.RequireApproval {
		if username != "" && username != userID {
			fmt.Printf("📝 Send DM to %s (%s):\n", username, userID)
		} else {
			fmt.Printf("📝 Send DM to user %s:\n", userID)
		}
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
				"action":    "send_dm",
				"cancelled": true,
			}, pretty)
		}
	}

	// Send DM
	msg, err := client.SendDirectMessage(userID, content)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(msg, pretty)
}

// isNumeric checks if a string contains only digits
func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func runDMHistory(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	limit, _ := cmd.Flags().GetInt("limit")
	userIdentifier := args[0]

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

	// Resolve username to user ID if needed
	userID := userIdentifier
	if !isNumeric(userIdentifier) {
		user, err := client.FindUserByUsername(userIdentifier)
		if err != nil {
			return output.PrintError(err, pretty)
		}
		userID = user.ID
	}

	// Get DM history
	messages, err := client.GetDMHistory(userID, limit)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"messages": messages,
		"count":    len(messages),
	}, pretty)
}

func runDMList(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	limit, _ := cmd.Flags().GetInt("limit")
	activeOnly, _ := cmd.Flags().GetBool("active-only")

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

	// Get DM channels
	dmChannels, err := client.ListDMChannels(limit, activeOnly)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"dm_channels": dmChannels,
		"count":       len(dmChannels),
	}, pretty)
}
