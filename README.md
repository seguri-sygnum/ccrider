# ccrider

Search, browse, and resume your Claude Code sessions. Fast.

## Why ccrider?

You've got months of Claude Code sessions sitting in `~/.claude/projects/`. Finding that conversation where you fixed the authentication bug? Good luck grepping through nested JSON files.

ccrider solves this with a TUI browser, CLI search, and an MCP server so Claude can search your past sessions too.

```bash
# Import your sessions once
ccrider sync

# Launch the TUI - browse, search, resume
ccrider tui

# Or search from command line
ccrider search "authentication bug"
```

Stay in your terminal. Find any conversation. Resume where you left off.

**Installation:**

```bash
# Homebrew (recommended)
brew install neilberkman/tap/ccrider

# Or from source
git clone https://github.com/neilberkman/ccrider.git
cd ccrider
go build -o ccrider cmd/ccrider/main.go
sudo mv ccrider /usr/local/bin/
```

_"Vibe code like ~a king~ The King!"_

https://github.com/user-attachments/assets/5b008290-076e-4323-a775-f27f704b1ff2

## Core Features

### 1. Interactive TUI Browser

```bash
ccrider tui
```

Browse your sessions with a polished terminal UI:

- **Arrow keys** to navigate
- **Enter** to view full conversation
- **o** to open session in new terminal tab (auto-detects Ghostty, iTerm, Terminal.app)
- **/** to search across all messages
- **p** to toggle project filter (show only current directory)
- **?** for help

Sessions matching your current directory are highlighted in light green - instantly see which sessions are relevant to your current work.

### 2. Full-Text Search

```bash
ccrider search "postgres migration"
ccrider search "error handling" --project ~/code/myapp
ccrider search "authentication" --after 2024-01-01
```

Powered by SQLite FTS5 - search message content, filter by project or date, get results instantly.

### 3. Resume Sessions

Press **r** in the TUI or use the CLI:

```bash
ccrider resume <session-id>
```

Launches `claude --resume` in the right directory with the right session. Just works.

### 4. Incremental Sync

```bash
ccrider sync       # Import all new sessions
ccrider sync --full  # Re-import everything
```

Detects ongoing sessions and imports new messages without re-processing everything.

---

## MCP Server

ccrider includes a built-in MCP (Model Context Protocol) server that gives Claude access to your session history.

Ask Claude to search your past conversations while working on new problems:

- "Find sessions where I worked on authentication"
- "Show me my most recent Elixir sessions"
- "What was I working on last week in the billing project?"
- "Search my sessions for postgres migration issues"

### Setup

**Claude Code:**

```bash
# Install for all your projects (recommended)
claude mcp add --scope user ccrider $(which ccrider) serve-mcp

# Or for current project only
claude mcp add ccrider $(which ccrider) serve-mcp
```

**Claude Desktop:**

Add to your config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "ccrider": {
      "command": "ccrider",
      "args": ["serve-mcp"]
    }
  }
}
```

### Available Tools

- **search_sessions** - Full-text search across all session content with date/project filters
- **list_recent_sessions** - Get recent sessions, optionally filtered by project
- **get_session_detail** - Retrieve session info with first/last messages and optional search
- **get_session_messages** - Get messages from a session (supports tail mode, context around search matches)

The MCP server provides read-only access to your session database. Your conversations stay local.

---

## Configuration

ccrider looks for config in `~/.config/ccrider/`:

```toml
# config.toml - pass additional flags to claude --resume
claude_flags = ["--dangerously-skip-permissions"]
```

```txt
# terminal_command.txt - custom command for 'o' key
# Available placeholders: {cwd}, {command}
wezterm cli spawn --cwd {cwd} -- {command}
```

```txt
# resume_prompt.txt - customize the prompt sent when resuming sessions
```

See [CONFIGURATION.md](docs/CONFIGURATION.md) for full details.

---

## Architecture

Built with strict core/interface separation following [Saša Jurić's principles](https://www.theerlangelist.com/article/phoenix_is_modular):

- **Core** (`pkg/`, `internal/core/`): Pure business logic - parsing, database, search
- **Interface** (`internal/interface/`, `cmd/`): Thin wrappers - CLI, TUI, MCP server

Uses proven technologies:

- **Go** for performance and single-binary distribution
- **SQLite with FTS5** for fast full-text search
- **Bubbletea** for polished terminal UI
- **MCP** for Claude integration

### Why This Matters

Other Claude Code session tools are broken:

- Incomplete schema support (can't parse all message types)
- Broken builds and abandoned dependencies
- No real search (just grep)
- Can't actually resume sessions

ccrider fixes this with:

- ✅ 100% schema coverage - parses all message types correctly
- ✅ SQLite FTS5 search - fast, powerful full-text search
- ✅ Single binary - no npm, no pip, no dependencies
- ✅ Native resume - one keystroke to resume sessions
- ✅ Incremental sync - detects new messages in ongoing sessions

---

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

### Project Structure

```
cmd/ccrider/          # CLI entry point + MCP server
internal/
  core/               # Business logic (no UI concerns)
    db/               # Database operations
    importer/         # Session import/sync
    search/           # Full-text search
    session/          # Session launch logic
  interface/          # Thin UI wrappers
    cli/              # Command handlers
    tui/              # Terminal UI (bubbletea)
pkg/ccsessions/       # Session file parser (public API)
```

### Quick Build

```bash
go build -o ccrider cmd/ccrider/main.go
./ccrider sync
./ccrider tui
```

---

## Documentation

- [Configuration Guide](docs/CONFIGURATION.md)
- [Resume Prompts](docs/RESUME_PROMPT.md)
- [Design Document](docs/plans/2025-11-08-ccrider-design.md)
- [Schema Documentation](research/schema.md)
- [Requirements](research/requirements.md)

## License

MIT
