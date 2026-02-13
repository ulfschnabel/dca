package main

import (
	"github.com/spf13/cobra"
	"github.com/ulfschnabel/dca/internal/config"
	"github.com/ulfschnabel/dca/internal/discord"
	"github.com/ulfschnabel/dca/internal/output"
	"fmt"
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Activity operations",
	Long:  "View recent activity across all servers and DMs",
}

var activityRecentCmd = &cobra.Command{
	Use:   "recent",
	Short: "Show recent activity",
	Long:  "Show recent messages across all servers and DMs, sorted by timestamp",
	RunE:  runActivityRecent,
}

func init() {
	rootCmd.AddCommand(activityCmd)
	activityCmd.AddCommand(activityRecentCmd)

	activityRecentCmd.Flags().Int("limit", 15, "Total messages to show")
	activityRecentCmd.Flags().String("type", "all", "Filter by type: all, dm, server")
}

func runActivityRecent(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	limit, _ := cmd.Flags().GetInt("limit")
	filterType, _ := cmd.Flags().GetString("type")

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

	// Get recent activity
	activity, err := client.GetRecentActivity(limit, filterType)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"activity": activity,
		"count":    len(activity),
	}, pretty)
}
