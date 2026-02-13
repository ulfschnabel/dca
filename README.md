# dca — Discord CLI for Agents

A unified CLI tool for interacting with Discord using your own account, designed for AI agents and personal automation.

## ⚠️ Important Disclaimer

**This tool uses Discord user tokens for automation, which is against Discord's Terms of Service.** This is for **personal use only** - similar to how slka uses user tokens for Slack automation.

- ✅ Use for personal automation
- ✅ Use to help AI agents interact with your Discord
- ❌ Do NOT use for spam, abuse, or commercial purposes
- ❌ Do NOT share your user token
- ⚠️ Use at your own risk

## Features

- **Acts as YOU** — Messages appear from your account, not a bot
- **JSON output** — All commands output JSON for easy parsing by LLMs
- **Human approval** — Write operations require confirmation
- **Token-efficient** — Optimized for AI agent token usage
- **Text-only** — Focused on messages and channels (no voice)

## Status

✅ **Phase 2 Complete** - Feature-complete for daily Discord use
- Activity: Recent messages across all DMs and servers
- Servers: List, info
- Channels: List, history
- Messages: Send, reply, edit, delete (with approval)
- DMs: List conversations, send, history
- Reactions: Add, remove

## Quick Start

### Installation

```bash
brew install ulfschnabel/tap/dca
```

Or download from [releases](https://github.com/ulfschnabel/dca/releases).

### Setup

```bash
# Initialize configuration
dca config init

# Follow prompts to enter your user token
```

### Getting Your User Token

1. Open Discord in browser (web.discord.com)
2. Press F12 to open Developer Tools
3. Go to Network tab
4. Send any message in Discord
5. Look for API requests (like `/messages`)
6. Find the `authorization` header in Request Headers
7. Copy that value - that's your user token

**Keep it secret!** This token gives full access to your Discord account.

## Commands

### Activity Overview (Start Here!) 🔥
```bash
dca activity recent --limit 15                 # See what's new everywhere
dca activity recent --type dm                  # Only DMs
dca activity recent --type server              # Only servers
```

### Direct Messages
```bash
dca dm list --limit 20                         # List DM conversations (sorted by activity)
dca dm history <username> --limit 10           # Get DM history
dca dm send <username> "text"                  # Send DM (finds user automatically)
```

### Messages
```bash
dca message send <channel-id> "text" --dry-run # Preview before sending
dca message send <channel-id> "text"           # Send message
dca message reply <channel-id> <msg-id> "text" # Reply to message
dca message edit <channel-id> <msg-id> "new"   # Edit your message
dca message delete <channel-id> <msg-id>       # Delete your message
```

### Reactions
```bash
dca reaction add <channel-id> <msg-id> 👍      # Add reaction
dca reaction remove <channel-id> <msg-id> 👍   # Remove reaction
```

### Servers & Channels
```bash
dca servers list                               # List your servers
dca servers info <server-id>                   # Server details
dca channels list <server-id>                  # List channels
dca channels history <channel-id> --limit 10   # Get messages
```

## For AI Agents

All commands return JSON:

```json
{
  "ok": true,
  "data": {
    // Command-specific data
  }
}
```

Perfect for LLM parsing and automation.

### Token Efficiency

dca is optimized for token-efficient AI agent use:

- **Small defaults**: `--limit` defaults to 10-20 (not 50-100)
- **Sorted by recency**: Latest activity first
- **Context included**: Server/channel names in activity feed
- **Filtering**: `--type`, `--active-only` to reduce noise

### Common Workflows

**Catch up on Discord:**
```bash
dca activity recent --limit 20      # What's new?
dca dm list --limit 10              # Who messaged me?
dca dm history alice --limit 5      # Read conversation
dca dm send alice "Got it, thanks!" # Respond
```

**Monitor specific channel:**
```bash
dca channels history <channel-id> --limit 5
dca message send <channel-id> "Update: ..."
```

**Quick reactions:**
```bash
dca activity recent --limit 10      # See latest
dca reaction add <ch-id> <msg-id> 👍 # React
```

## Development

```bash
# Build
go build ./cmd/dca

# Test
./dca config init
./dca servers list
```

## Why Not a Bot?

**User token** (this tool):
- Messages appear as YOU
- Access your servers automatically
- Can read your DMs
- Personal automation

**Bot token**:
- Messages tagged as "BOT"
- Must invite to servers
- Can't read your DMs
- More setup

This tool is for personal automation, like having an assistant check Discord for you.

## License

MIT

## Disclaimer

This project is not affiliated with Discord. Use responsibly and at your own risk.
