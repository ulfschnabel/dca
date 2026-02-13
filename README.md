# dca — Discord CLI for Agentic Workflows

A unified CLI tool for interacting with Discord, designed for AI agents and automation workflows.

## Features

- **Unified binary** — Single `dca` command with both read and write operations
- **JSON output** — All commands output JSON for easy parsing by LLMs and scripts
- **Human approval mode** — Write operations can require explicit human confirmation
- **Token-efficient** — Optimized for AI agent token usage
- **Bot-based** — Uses Discord bot tokens for reliable API access

## Status

🚧 **Under Development** - Initial implementation in progress

## Planned Commands

### Servers (Guilds)
```bash
dca servers list                    # List all servers bot has access to
dca servers info <server-id>        # Get server information
```

### Channels
```bash
dca channels list <server-id>       # List channels in a server
dca channels info <channel-id>      # Get channel information
dca channels history <channel-id>   # Get message history
```

### Messages
```bash
dca message send <channel-id> "text"           # Send message
dca message reply <channel-id> <msg-id> "text" # Reply to message
dca message edit <channel-id> <msg-id> "text"  # Edit message
dca message delete <channel-id> <msg-id>       # Delete message
```

### Direct Messages
```bash
dca dm list                         # List DM channels
dca dm send <user-id> "text"        # Send DM
dca dm history <user-id>            # Get DM history
```

### Reactions
```bash
dca reaction list <channel-id> <msg-id>        # List reactions
dca reaction add <channel-id> <msg-id> :emoji: # Add reaction
dca reaction remove <channel-id> <msg-id> :emoji: # Remove reaction
```

### Users
```bash
dca users lookup <username>         # Find user by username
dca users info <user-id>            # Get user information
```

## Installation

Coming soon - will be available via Homebrew and direct download.

## Configuration

Config will be stored in `~/.config/dca/config.json`:

```json
{
  "bot_token": "your-bot-token-here",
  "require_approval": true
}
```

## For AI Agents

All commands output JSON for easy parsing. Token-efficient design optimized for LLM usage.

## Development

```bash
# Build
go build ./cmd/dca

# Run
./dca --help
```

## License

MIT
