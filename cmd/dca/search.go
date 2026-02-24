package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ulfschnabel/dca/internal/config"
	"github.com/ulfschnabel/dca/internal/discord"
	"github.com/ulfschnabel/dca/internal/output"
)

var searchCmd = &cobra.Command{
	Use:   "search <server-id> <query>",
	Short: "Search messages in a server",
	Long: `Search for messages in a Discord server by keyword.

Results are paginated in groups of 25. Use --offset to get more pages.

Examples:
  dca search 123456789 "error log"
  dca search 123456789 "deployment" --channel-id 987654321
  dca search 123456789 "bug" --author-id 111222333 --sort-by timestamp`,
	Args: cobra.ExactArgs(2),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().String("channel-id", "", "Filter results to a specific channel")
	searchCmd.Flags().String("author-id", "", "Filter results by author ID")
	searchCmd.Flags().String("has", "", "Filter by attachment type: link, embed, file, video, image, sound")
	searchCmd.Flags().Int("offset", 0, "Pagination offset (multiples of 25)")
	searchCmd.Flags().String("sort-by", "", "Sort by: relevance or timestamp")
	searchCmd.Flags().String("sort-order", "", "Sort order: asc or desc")
}

func runSearch(cmd *cobra.Command, args []string) error {
	pretty, _ := cmd.Flags().GetBool("output-pretty")
	serverID := args[0]
	query := args[1]

	cfg, err := config.Load(cfgFile)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		token = cfg.UserToken
	}

	if token == "" {
		return output.PrintError(fmt.Errorf("no token configured"), pretty)
	}

	client, err := discord.New(token)
	if err != nil {
		return output.PrintError(err, pretty)
	}
	defer client.Close()

	channelID, _ := cmd.Flags().GetString("channel-id")
	authorID, _ := cmd.Flags().GetString("author-id")
	has, _ := cmd.Flags().GetString("has")
	offset, _ := cmd.Flags().GetInt("offset")
	sortBy, _ := cmd.Flags().GetString("sort-by")
	sortOrder, _ := cmd.Flags().GetString("sort-order")

	opts := discord.SearchOptions{
		Content:   query,
		AuthorID:  authorID,
		ChannelID: channelID,
		Has:       has,
		Offset:    offset,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}

	result, err := client.SearchGuildMessages(serverID, opts)
	if err != nil {
		return output.PrintError(err, pretty)
	}

	return output.PrintSuccess(map[string]interface{}{
		"messages":      result.Messages,
		"count":         len(result.Messages),
		"total_results": result.TotalResults,
		"offset":        offset,
	}, pretty)
}
