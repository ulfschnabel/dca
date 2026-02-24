package discord

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Client wraps the Discord session
type Client struct {
	session *discordgo.Session
}

// New creates a new Discord client with a user token
func New(token string) (*Client, error) {
	// Create session with user token (no "Bot " prefix for user tokens)
	session, err := discordgo.New(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &Client{session: session}, nil
}

// Close closes the Discord session
func (c *Client) Close() error {
	return c.session.Close()
}

// Guild represents a Discord server
type Guild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MemberCount int    `json:"member_count"`
	OwnerID     string `json:"owner_id"`
}

// ListGuilds returns all guilds the bot has access to
func (c *Client) ListGuilds() ([]*Guild, error) {
	guilds, err := c.session.UserGuilds(100, "", "", false)
	if err != nil {
		return nil, fmt.Errorf("failed to list guilds: %w", err)
	}

	result := make([]*Guild, 0, len(guilds))
	for _, g := range guilds {
		result = append(result, &Guild{
			ID:   g.ID,
			Name: g.Name,
			// UserGuild doesn't have Description, OwnerID, or MemberCount
			// Use GetGuild() for detailed information
		})
	}

	return result, nil
}

// GetGuild gets detailed information about a guild
func (c *Client) GetGuild(guildID string) (*Guild, error) {
	g, err := c.session.Guild(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild: %w", err)
	}

	return &Guild{
		ID:          g.ID,
		Name:        g.Name,
		Description: g.Description,
		MemberCount: g.ApproximateMemberCount,
		OwnerID:     g.OwnerID,
	}, nil
}

// Channel represents a Discord channel
type Channel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Topic    string `json:"topic,omitempty"`
	ParentID string `json:"parent_id,omitempty"`
}

// ForumThread represents a thread in a forum channel
type ForumThread struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	MessageCount  int     `json:"message_count"`
	CreatedAt     string  `json:"created_at"`
	LastMessageID string  `json:"last_message_id,omitempty"`
	Archived      bool    `json:"archived"`
	Author        *Author `json:"author,omitempty"`
}

// channelTypeToString converts a ChannelType to a readable string
func channelTypeToString(t discordgo.ChannelType) string {
	switch t {
	case discordgo.ChannelTypeGuildText:
		return "text"
	case discordgo.ChannelTypeDM:
		return "dm"
	case discordgo.ChannelTypeGuildVoice:
		return "voice"
	case discordgo.ChannelTypeGroupDM:
		return "group_dm"
	case discordgo.ChannelTypeGuildCategory:
		return "category"
	case discordgo.ChannelTypeGuildNews:
		return "news"
	case discordgo.ChannelTypeGuildStore:
		return "store"
	case discordgo.ChannelTypeGuildNewsThread:
		return "news_thread"
	case discordgo.ChannelTypeGuildPublicThread:
		return "public_thread"
	case discordgo.ChannelTypeGuildPrivateThread:
		return "private_thread"
	case discordgo.ChannelTypeGuildStageVoice:
		return "stage"
	case discordgo.ChannelTypeGuildForum:
		return "forum"
	default:
		return fmt.Sprintf("unknown_%d", t)
	}
}

// ListChannels returns all channels in a guild
func (c *Client) ListChannels(guildID string) ([]*Channel, error) {
	channels, err := c.session.GuildChannels(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	result := make([]*Channel, 0, len(channels))
	for _, ch := range channels {
		result = append(result, &Channel{
			ID:       ch.ID,
			Name:     ch.Name,
			Type:     channelTypeToString(ch.Type),
			Topic:    ch.Topic,
			ParentID: ch.ParentID,
		})
	}

	return result, nil
}

// Message represents a Discord message
type Message struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Author    Author `json:"author"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// Author represents a message author
type Author struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bot      bool   `json:"bot"`
}

// GetMessages retrieves messages from a channel
func (c *Client) GetMessages(channelID string, limit int) ([]*Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	msgs, err := c.session.ChannelMessages(channelID, limit, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	result := make([]*Message, 0, len(msgs))
	for _, m := range msgs {
		result = append(result, &Message{
			ID:        m.ID,
			ChannelID: m.ChannelID,
			Author: Author{
				ID:       m.Author.ID,
				Username: m.Author.Username,
				Bot:      m.Author.Bot,
			},
			Content:   m.Content,
			Timestamp: m.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return result, nil
}

// SendMessage sends a message to a channel
func (c *Client) SendMessage(channelID, content string) (*Message, error) {
	msg, err := c.session.ChannelMessageSend(channelID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return &Message{
		ID:        msg.ID,
		ChannelID: msg.ChannelID,
		Author: Author{
			ID:       msg.Author.ID,
			Username: msg.Author.Username,
			Bot:      msg.Author.Bot,
		},
		Content:   msg.Content,
		Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ReplyToMessage replies to a specific message
func (c *Client) ReplyToMessage(channelID, messageID, content string) (*Message, error) {
	msg, err := c.session.ChannelMessageSendReply(channelID, content, &discordgo.MessageReference{
		MessageID: messageID,
		ChannelID: channelID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to reply to message: %w", err)
	}

	return &Message{
		ID:        msg.ID,
		ChannelID: msg.ChannelID,
		Author: Author{
			ID:       msg.Author.ID,
			Username: msg.Author.Username,
			Bot:      msg.Author.Bot,
		},
		Content:   msg.Content,
		Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// SendDirectMessage sends a DM to a user
func (c *Client) SendDirectMessage(userID, content string) (*Message, error) {
	// Create or get DM channel with user
	channel, err := c.session.UserChannelCreate(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create DM channel: %w", err)
	}

	// Send message to DM channel
	msg, err := c.session.ChannelMessageSend(channel.ID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to send DM: %w", err)
	}

	return &Message{
		ID:        msg.ID,
		ChannelID: msg.ChannelID,
		Author: Author{
			ID:       msg.Author.ID,
			Username: msg.Author.Username,
			Bot:      msg.Author.Bot,
		},
		Content:   msg.Content,
		Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetDMHistory gets message history from DMs with a user
func (c *Client) GetDMHistory(userID string, limit int) ([]*Message, error) {
	// Create or get DM channel with user
	channel, err := c.session.UserChannelCreate(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create DM channel: %w", err)
	}

	// Get messages
	return c.GetMessages(channel.ID, limit)
}

// DMChannel represents a DM conversation
type DMChannel struct {
	ChannelID   string   `json:"channel_id"`
	User        Author   `json:"user"`
	LastMessage *Message `json:"last_message,omitempty"`
}

// ActivityMessage represents a message with full context
type ActivityMessage struct {
	Message
	Type        string  `json:"type"` // "dm" or "server"
	ServerName  string  `json:"server_name,omitempty"`
	ServerID    string  `json:"server_id,omitempty"`
	ChannelName string  `json:"channel_name,omitempty"`
	DMUser      *Author `json:"dm_user,omitempty"`
}

// EditMessage edits a message
func (c *Client) EditMessage(channelID, messageID, newContent string) (*Message, error) {
	msg, err := c.session.ChannelMessageEdit(channelID, messageID, newContent)
	if err != nil {
		return nil, fmt.Errorf("failed to edit message: %w", err)
	}

	return &Message{
		ID:        msg.ID,
		ChannelID: msg.ChannelID,
		Author: Author{
			ID:       msg.Author.ID,
			Username: msg.Author.Username,
			Bot:      msg.Author.Bot,
		},
		Content:   msg.Content,
		Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteMessage deletes a message
func (c *Client) DeleteMessage(channelID, messageID string) error {
	err := c.session.ChannelMessageDelete(channelID, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

// GetMessage retrieves a specific message
func (c *Client) GetMessage(channelID, messageID string) (*Message, error) {
	msg, err := c.session.ChannelMessage(channelID, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &Message{
		ID:        msg.ID,
		ChannelID: msg.ChannelID,
		Author: Author{
			ID:       msg.Author.ID,
			Username: msg.Author.Username,
			Bot:      msg.Author.Bot,
		},
		Content:   msg.Content,
		Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// AddReaction adds a reaction to a message
func (c *Client) AddReaction(channelID, messageID, emoji string) error {
	err := c.session.MessageReactionAdd(channelID, messageID, emoji)
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}
	return nil
}

// RemoveReaction removes your reaction from a message
func (c *Client) RemoveReaction(channelID, messageID, emoji string) error {
	err := c.session.MessageReactionRemove(channelID, messageID, emoji, "@me")
	if err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}
	return nil
}

// GetRecentActivity gets recent messages across all DMs and servers
func (c *Client) GetRecentActivity(limit int, filterType string) ([]*ActivityMessage, error) {
	var allMessages []*ActivityMessage

	// Get DM messages if requested
	if filterType == "all" || filterType == "dm" {
		dmMessages, err := c.getRecentDMMessages(5) // Get recent from each DM
		if err == nil {
			allMessages = append(allMessages, dmMessages...)
		}
	}

	// Get server messages if requested
	if filterType == "all" || filterType == "server" {
		serverMessages, err := c.getRecentServerMessages(5) // Get recent from each channel
		if err == nil {
			allMessages = append(allMessages, serverMessages...)
		}
	}

	// Sort by timestamp (newest first)
	for i := 0; i < len(allMessages)-1; i++ {
		for j := 0; j < len(allMessages)-i-1; j++ {
			if allMessages[j].Timestamp < allMessages[j+1].Timestamp {
				allMessages[j], allMessages[j+1] = allMessages[j+1], allMessages[j]
			}
		}
	}

	// Apply limit
	if limit > 0 && len(allMessages) > limit {
		allMessages = allMessages[:limit]
	}

	return allMessages, nil
}

func (c *Client) getRecentDMMessages(perChannel int) ([]*ActivityMessage, error) {
	var messages []*ActivityMessage

	// Get DM channels
	var channels []*discordgo.Channel
	body, err := c.session.RequestWithBucketID("GET", discordgo.EndpointUserChannels("@me"), nil, discordgo.EndpointUserChannels(""))
	if err != nil {
		return nil, err
	}

	if err = discordgo.Unmarshal(body, &channels); err != nil {
		return nil, err
	}

	// Get messages from each DM
	for _, ch := range channels {
		if ch.Type != discordgo.ChannelTypeDM {
			continue
		}

		msgs, err := c.session.ChannelMessages(ch.ID, perChannel, "", "", "")
		if err != nil || len(msgs) == 0 {
			continue
		}

		// Get DM user
		var dmUser *Author
		if len(ch.Recipients) > 0 {
			dmUser = &Author{
				ID:       ch.Recipients[0].ID,
				Username: ch.Recipients[0].Username,
				Bot:      ch.Recipients[0].Bot,
			}
		}

		for _, msg := range msgs {
			messages = append(messages, &ActivityMessage{
				Message: Message{
					ID:        msg.ID,
					ChannelID: msg.ChannelID,
					Author: Author{
						ID:       msg.Author.ID,
						Username: msg.Author.Username,
						Bot:      msg.Author.Bot,
					},
					Content:   msg.Content,
					Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
				},
				Type:   "dm",
				DMUser: dmUser,
			})
		}
	}

	return messages, nil
}

// ListForumThreads lists threads in a forum channel (archived and active)
func (c *Client) ListForumThreads(channelID string, limit int, activeOnly bool) ([]*ForumThread, error) {
	// Try to get archived public threads (works with user tokens)
	// This is a workaround since active threads endpoint is bot-only
	endpoint := fmt.Sprintf("%s/channels/%s/threads/archived/public", discordgo.EndpointAPI, channelID)

	type ThreadResponse struct {
		Threads []*discordgo.Channel `json:"threads"`
		HasMore bool                 `json:"has_more"`
	}

	var response ThreadResponse
	body, err := c.session.RequestWithBucketID("GET", endpoint, nil, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get forum threads (note: user tokens can only see archived public threads): %w", err)
	}

	if err = discordgo.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse forum threads: %w", err)
	}

	result := make([]*ForumThread, 0)
	for _, thread := range response.Threads {
		// Skip archived if activeOnly
		if activeOnly && thread.ThreadMetadata != nil && thread.ThreadMetadata.Archived {
			continue
		}

		forumThread := &ForumThread{
			ID:            thread.ID,
			Name:          thread.Name,
			MessageCount:  thread.MessageCount,
			LastMessageID: thread.LastMessageID,
			Archived:      thread.ThreadMetadata != nil && thread.ThreadMetadata.Archived,
		}

		result = append(result, forumThread)

		if limit > 0 && len(result) >= limit {
			break
		}
	}

	return result, nil
}

// GetThreadMessages gets messages from a specific thread
func (c *Client) GetThreadMessages(threadID string, limit int) ([]*Message, error) {
	// Threads are just channels, so we can use the regular GetMessages
	return c.GetMessages(threadID, limit)
}

func (c *Client) getRecentServerMessages(perChannel int) ([]*ActivityMessage, error) {
	var messages []*ActivityMessage

	// Get guilds
	guilds, err := c.session.UserGuilds(100, "", "", false)
	if err != nil {
		return nil, err
	}

	// Get messages from each guild's channels (limit to avoid rate limits)
	for _, guild := range guilds {
		channels, err := c.session.GuildChannels(guild.ID)
		if err != nil {
			continue
		}

		// Only check first 3 text channels per guild for efficiency
		checked := 0
		for _, ch := range channels {
			if ch.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			if checked >= 3 {
				break
			}
			checked++

			msgs, err := c.session.ChannelMessages(ch.ID, perChannel, "", "", "")
			if err != nil {
				continue
			}

			for _, msg := range msgs {
				messages = append(messages, &ActivityMessage{
					Message: Message{
						ID:        msg.ID,
						ChannelID: msg.ChannelID,
						Author: Author{
							ID:       msg.Author.ID,
							Username: msg.Author.Username,
							Bot:      msg.Author.Bot,
						},
						Content:   msg.Content,
						Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
					},
					Type:        "server",
					ServerName:  guild.Name,
					ServerID:    guild.ID,
					ChannelName: ch.Name,
				})
			}
		}
	}

	return messages, nil
}

// ListDMChannels returns all DM channels sorted by recent activity
func (c *Client) ListDMChannels(limit int, activeOnly bool) ([]*DMChannel, error) {
	// Get user's DM channels
	var channels []*discordgo.Channel
	body, err := c.session.RequestWithBucketID("GET", discordgo.EndpointUserChannels("@me"), nil, discordgo.EndpointUserChannels(""))
	if err != nil {
		return nil, fmt.Errorf("failed to get DM channels: %w", err)
	}

	if err = discordgo.Unmarshal(body, &channels); err != nil {
		return nil, fmt.Errorf("failed to parse DM channels: %w", err)
	}

	result := make([]*DMChannel, 0)
	for _, ch := range channels {
		if ch.Type != discordgo.ChannelTypeDM {
			continue
		}

		// Get last message
		msgs, err := c.session.ChannelMessages(ch.ID, 1, "", "", "")
		if err != nil || len(msgs) == 0 {
			if activeOnly {
				continue // Skip DMs with no messages if activeOnly
			}
			// Include DM with no last message
			if len(ch.Recipients) > 0 {
				result = append(result, &DMChannel{
					ChannelID: ch.ID,
					User: Author{
						ID:       ch.Recipients[0].ID,
						Username: ch.Recipients[0].Username,
						Bot:      ch.Recipients[0].Bot,
					},
				})
			}
			continue
		}

		lastMsg := msgs[0]
		dmChannel := &DMChannel{
			ChannelID: ch.ID,
			LastMessage: &Message{
				ID:        lastMsg.ID,
				ChannelID: lastMsg.ChannelID,
				Author: Author{
					ID:       lastMsg.Author.ID,
					Username: lastMsg.Author.Username,
					Bot:      lastMsg.Author.Bot,
				},
				Content:   lastMsg.Content,
				Timestamp: lastMsg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			},
		}

		// Set user (recipient)
		if len(ch.Recipients) > 0 {
			dmChannel.User = Author{
				ID:       ch.Recipients[0].ID,
				Username: ch.Recipients[0].Username,
				Bot:      ch.Recipients[0].Bot,
			}
		}

		result = append(result, dmChannel)
	}

	// Sort by last message timestamp (newest first)
	// Simple bubble sort since we expect small lists
	for i := 0; i < len(result)-1; i++ {
		for j := 0; j < len(result)-i-1; j++ {
			if result[j].LastMessage == nil {
				continue
			}
			if result[j+1].LastMessage == nil {
				// Swap if next has no message
				result[j], result[j+1] = result[j+1], result[j]
				continue
			}
			if result[j].LastMessage.Timestamp < result[j+1].LastMessage.Timestamp {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}

	// Apply limit
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// FindUserByUsername searches for a user by username across DMs and guilds
func (c *Client) FindUserByUsername(username string) (*Author, error) {
	// First check DM channels (most likely location)
	fmt.Printf("   Searching DM channels...\n")

	// Get user's DM channels using REST API
	var dmChannels []*discordgo.Channel
	body, err := c.session.RequestWithBucketID("GET", discordgo.EndpointUserChannels("@me"), nil, discordgo.EndpointUserChannels(""))
	if err == nil {
		err = discordgo.Unmarshal(body, &dmChannels)
		if err == nil {
			for _, channel := range dmChannels {
				if channel.Type != discordgo.ChannelTypeDM {
					continue
				}

				// Get recent messages from DM
				msgs, err := c.session.ChannelMessages(channel.ID, 20, "", "", "")
				if err != nil {
					continue
				}

				// Look for username match
				for _, msg := range msgs {
					if strings.EqualFold(msg.Author.Username, username) {
						return &Author{
							ID:       msg.Author.ID,
							Username: msg.Author.Username,
							Bot:      msg.Author.Bot,
						}, nil
					}
				}
			}
		}
	}

	// If not found in DMs, search guilds
	guilds, err := c.session.UserGuilds(100, "", "", false)
	if err != nil {
		return nil, fmt.Errorf("failed to list guilds: %w", err)
	}

	fmt.Printf("   Searching %d servers...\n", len(guilds))

	// Search through each guild's channels for messages from this user
	for i, guild := range guilds {
		fmt.Printf("   [%d/%d] %s\n", i+1, len(guilds), guild.Name)

		channels, err := c.session.GuildChannels(guild.ID)
		if err != nil {
			continue // Skip inaccessible guilds
		}

		// Check text channels only, limit to first 5 channels per guild for speed
		checked := 0
		for _, channel := range channels {
			if channel.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			if checked >= 5 {
				break // Limit search depth
			}
			checked++

			// Get recent messages
			msgs, err := c.session.ChannelMessages(channel.ID, 20, "", "", "")
			if err != nil {
				continue // Skip inaccessible channels
			}

			// Look for username match
			for _, msg := range msgs {
				if strings.EqualFold(msg.Author.Username, username) {
					return &Author{
						ID:       msg.Author.ID,
						Username: msg.Author.Username,
						Bot:      msg.Author.Bot,
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("user %q not found in recent messages", username)
}

// SearchOptions configures a Discord message search query
type SearchOptions struct {
	Content   string
	AuthorID  string
	ChannelID string
	Has       string
	Offset    int
	SortBy    string
	SortOrder string
}

// SearchResult holds the parsed search response
type SearchResult struct {
	TotalResults int              `json:"total_results"`
	Messages     []*SearchMessage `json:"messages"`
}

// SearchMessage represents a single search hit
type SearchMessage struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Author    Author `json:"author"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// SearchGuildMessages searches for messages in a guild using Discord's search API
func (c *Client) SearchGuildMessages(guildID string, opts SearchOptions) (*SearchResult, error) {
	params := buildSearchParams(opts)
	baseEndpoint := fmt.Sprintf("%sguilds/%s/messages/search", discordgo.EndpointAPI, guildID)
	fullEndpoint := baseEndpoint + "?" + params.Encode()

	body, err := c.session.RequestWithBucketID("GET", fullEndpoint, nil, baseEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to search guild messages: %w", err)
	}

	return parseSearchResponse(body)
}

// buildSearchParams constructs URL query parameters from SearchOptions
func buildSearchParams(opts SearchOptions) url.Values {
	params := url.Values{}
	params.Set("content", opts.Content)
	if opts.AuthorID != "" {
		params.Set("author_id", opts.AuthorID)
	}
	if opts.ChannelID != "" {
		params.Set("channel_id", opts.ChannelID)
	}
	if opts.Has != "" {
		params.Set("has", opts.Has)
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.SortBy != "" {
		params.Set("sort_by", opts.SortBy)
	}
	if opts.SortOrder != "" {
		params.Set("sort_order", opts.SortOrder)
	}
	return params
}

// parseSearchResponse parses Discord's nested search response format.
// Discord returns messages as [[msg, msg], [msg, msg]] where each inner array
// is a group of context messages. The matched message has "hit": true.
func parseSearchResponse(body []byte) (*SearchResult, error) {
	var raw struct {
		TotalResults int                 `json:"total_results"`
		Messages     [][]json.RawMessage `json:"messages"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse search response: %w", err)
	}

	result := &SearchResult{
		TotalResults: raw.TotalResults,
		Messages:     make([]*SearchMessage, 0),
	}

	for _, group := range raw.Messages {
		for _, msgRaw := range group {
			var msg struct {
				ID        string `json:"id"`
				ChannelID string `json:"channel_id"`
				Content   string `json:"content"`
				Timestamp string `json:"timestamp"`
				Hit       bool   `json:"hit"`
				Author    struct {
					ID       string `json:"id"`
					Username string `json:"username"`
					Bot      bool   `json:"bot"`
				} `json:"author"`
			}
			if err := json.Unmarshal(msgRaw, &msg); err != nil {
				continue
			}
			if msg.Hit {
				result.Messages = append(result.Messages, &SearchMessage{
					ID:        msg.ID,
					ChannelID: msg.ChannelID,
					Author: Author{
						ID:       msg.Author.ID,
						Username: msg.Author.Username,
						Bot:      msg.Author.Bot,
					},
					Content:   msg.Content,
					Timestamp: msg.Timestamp,
				})
			}
		}
	}

	return result, nil
}
