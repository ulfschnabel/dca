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

🚧 **Under Development** - Phase 1 implementation in progress

## Quick Start

### Installation

Coming soon - will be available via Homebrew.

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

## Planned Commands

### Servers
```bash
dca servers list                    # List your servers
dca servers info <server-id>        # Server details
```

### Channels
```bash
dca channels list <server-id>       # List channels
dca channels history <channel-id>   # Get messages
```

### Messages
```bash
dca message send <channel-id> "text"           # Send message
dca message reply <channel-id> <msg-id> "text" # Reply
dca message edit <channel-id> <msg-id> "text"  # Edit your message
```

### Direct Messages
```bash
dca dm list                         # List DM channels
dca dm send <user-id> "text"        # Send DM
```

### Reactions
```bash
dca reaction add <channel-id> <msg-id> :emoji: # React
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
