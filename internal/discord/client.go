package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Client wraps the Discord session
type Client struct {
	session *discordgo.Session
}

// New creates a new Discord client
func New(token string) (*Client, error) {
	// Create session
	session, err := discordgo.New("Bot " + token)
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
			ID:          g.ID,
			Name:        g.Name,
			Description: g.Description,
			OwnerID:     g.OwnerID,
			// MemberCount needs approximate count from Guild object
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
