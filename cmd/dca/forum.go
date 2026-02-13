package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulfschnabel/dca/internal/config"
	"github.com/ulfschnabel/dca/internal/discord"
	"github.com/ulfschnabel/dca/internal/output"
)

var forumCmd = &cobra.Command{
	Use:   "forum",
	Short: "Forum channel operations",
	Long:  "List threads and read messages from forum channels",
}

var forumThreadsCmd = &cobra.Command{
	Use:   "threads <channel-id>",
	Short: "List forum threads",
	Long:  "List active threads in a forum channel",
	Args:  cobra.ExactArgs(1),
	RunE:  runForumThreads,
}

var forumMessagesCmd = &cobra.Command{
	Use:   "messages <thread-id>",
	Short: "Get messages from a thread",
	Long:  "Get messages from a specific forum thread",
	Args:  cobra.ExactArgs(1),
	RunE:  runForumMessages,
}

func init() {
	rootCmd.AddCommand(forumCmd)
	forumCmd.AddCommand(forumThreadsCmd)
	forumCmd.AddCommand(forumMessagesCmd)

	forumThreadsCmd.Flags().Int("limit", 20, "Number of threads to show")
	forumThreadsCmd.Flags().Bool("active-only", true, "Only show active (non-archived) threads")
	forumMessagesCmd.Flags().Int("limit", 10, "Number of messages to retrieve")
}

func runForumThreads(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	limit, _ := cmd.Flags().GetInt("limit")
	activeOnly, _ := cmd.Flags().GetBool("active-only")
	channelID := args[0]

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

	// List forum threads
	threads, err := client.ListForumThreads(channelID, limit, activeOnly)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"threads": threads,
		"count":   len(threads),
	}, pretty)
}

func runForumMessages(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	limit, _ := cmd.Flags().GetInt("limit")
	threadID := args[0]

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

	// Get thread messages
	messages, err := client.GetThreadMessages(threadID, limit)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"messages": messages,
		"count":    len(messages),
	}, pretty)
}
