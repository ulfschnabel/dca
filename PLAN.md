# dca - Discord CLI for Agents - Implementation Plan

## Overview

**Name**: dca (Discord CLI for Agents)
**Purpose**: Text-based Discord CLI for AI agents and automation
**Model**: Similar to slka - interact with Discord as you would normally
**Auth**: User token (personal automation, like slka's user tokens)

## Core Principles (from slka)

1. **JSON output** - All commands return structured JSON
2. **Token efficiency** - Optimized for LLM token usage
3. **Human approval** - Write operations require confirmation
4. **Dry-run mode** - Test commands before execution
5. **Single binary** - One tool, all operations
6. **Text-only** - No voice channel support
7. **User token auth** - Acts as you, not as a bot

## Important Note

⚠️ **Discord ToS Disclaimer**: Using user tokens for automation is against Discord's Terms of Service. This tool is for **personal automation only** - similar to how slka uses user tokens for Slack. Use at your own risk. Do not use for spam, abuse, or commercial purposes.

## Feature Roadmap

### Phase 1: MVP (First Release)

**Config Management**
- ✅ `dca config init` - Interactive setup
- ✅ `dca config show` - Display config (token masked)
- Store: bot token, approval setting

**Servers (Guilds)**
- `dca servers list` - List all servers bot has access to
- `dca servers info <server-id>` - Get server details
- Output: ID, name, member count, owner

**Channels**
- `dca channels list <server-id>` - List channels in server
- `dca channels info <channel-id>` - Get channel details
- `dca channels history <channel-id> --limit N` - Get messages

**Messages (Read)**
- `dca message get <channel-id> <message-id>` - Get specific message
- Included in channels history

**Messages (Write - with approval)**
- `dca message send <channel-id> "text"` - Send message
- `dca message reply <channel-id> <message-id> "text"` - Reply
- `dca message edit <channel-id> <message-id> "text"` - Edit own message
- All support `--dry-run`

### Phase 2: Enhanced Features

**Direct Messages**
- `dca dm list` - List DM channels
- `dca dm send <user-id> "text"` - Send DM
- `dca dm history <user-id>` - Get DM history

**Reactions**
- `dca reaction list <channel-id> <message-id>` - List reactions
- `dca reaction add <channel-id> <message-id> :emoji:` - Add reaction
- `dca reaction remove <channel-id> <message-id> :emoji:` - Remove

**Users**
- `dca users lookup <username>` - Find user by username
- `dca users info <user-id>` - Get user details

**Threads**
- `dca thread list <channel-id>` - List active threads
- `dca thread create <channel-id> <message-id> "name"` - Create thread
- `dca thread reply <thread-id> "text"` - Reply to thread

### Phase 3: Advanced Features

**Search & Filter**
- `dca channels list --filter <name>` - Filter channels by name
- `dca messages search <channel-id> "query"` - Search messages

**Bulk Operations**
- `dca message send-bulk <channel-ids...> "text"` - Send to multiple channels
- Rate limiting and batch handling

**Webhooks**
- `dca webhook create <channel-id> "name"` - Create webhook
- `dca webhook send <webhook-url> "text"` - Send via webhook

## Command Structure (Following slka Pattern)

```
dca
├── config
│   ├── init         # Interactive setup
│   └── show         # Display config
├── servers
│   ├── list         # List all servers
│   └── info         # Server details
├── channels
│   ├── list         # List channels in server
│   ├── info         # Channel details
│   └── history      # Get messages
├── message
│   ├── get          # Get specific message
│   ├── send         # Send message (approval)
│   ├── reply        # Reply to message (approval)
│   ├── edit         # Edit message (approval)
│   └── delete       # Delete message (approval)
├── dm
│   ├── list         # List DM channels
│   ├── send         # Send DM (approval)
│   └── history      # Get DM history
├── reaction
│   ├── list         # List reactions
│   ├── add          # Add reaction (approval)
│   └── remove       # Remove reaction (approval)
├── users
│   ├── lookup       # Find user by username
│   └── info         # User details
└── thread
    ├── list         # List threads
    ├── create       # Create thread (approval)
    └── reply        # Reply to thread (approval)
```

## JSON Output Format

All commands follow this structure:

```json
{
  "ok": true,
  "data": {
    // Command-specific data
  },
  "error": null  // Only present if ok: false
}
```

## Configuration File

`~/.config/dca/config.json`:

```json
{
  "user_token": "your-user-token-here",
  "require_approval": true
}
```

### How to Get Your User Token

1. Open Discord in your browser (web.discord.com)
2. Open Developer Tools (F12)
3. Go to Network tab
4. Send a message or reload
5. Look for any API request
6. In Request Headers, find `authorization` header
7. Copy that token (it's your user token)

**Security**: Keep this token private - it gives full access to your Discord account!

## Technical Stack

- **Language**: Go
- **CLI Framework**: Cobra
- **Discord Library**: discordgo
- **Output**: JSON
- **Distribution**: Homebrew + direct download

## User Token vs Bot Token

**Why user token?**
- Acts as YOU - no "bot" tag on messages
- Access to your servers automatically
- Can read DMs
- Personal automation (like slka)

**Why NOT bot token?**
- Messages show as from a bot
- Need to invite bot to servers
- Can't read your DMs
- More setup overhead

This tool is designed for **personal use**, not for running public bots.

## Development Phases

### Week 1: Core Infrastructure
- ✅ Project setup (dca)
- ✅ Config management
- ✅ Discord client wrapper
- Basic error handling
- JSON output utilities

### Week 2: Phase 1 Features
- Servers commands
- Channels commands
- Messages (read)
- Messages (write with approval)
- Testing with real Discord server

### Week 3: Phase 2 Features
- DMs
- Reactions
- Users
- Threads (basic)

### Week 4: Polish & Release
- Documentation
- Bot setup guide
- GoReleaser config
- Homebrew formula
- Release v0.1.0

## Testing Strategy

1. **Unit Tests**: Core logic, JSON parsing
2. **Integration Tests**: Discord API calls (with test bot)
3. **Manual Testing**: Real Discord server for agent scenarios

## Success Metrics

- Can interact with Discord via CLI as naturally as slka does with Slack
- AI agents can parse JSON output reliably
- Approval flow prevents accidental writes
- Token-efficient for LLM usage

## Open Questions

1. How to handle Discord's rate limits?
2. Should we support embed messages?
3. Do we need file upload support?
4. How to handle large servers (1000+ channels)?

## Next Steps

1. Finish Phase 1 implementation
2. Test with a real Discord server
3. Get feedback on command structure
4. Iterate based on actual use cases
