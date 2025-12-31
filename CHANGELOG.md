# Changelog

## [0.9.1] - 2024-12-30

### Fixed

- Phrase search now works in TUI (quotes were being stripped)
- Auto-balance quotes during live typing to prevent FTS5 errors
- Restored Elvis video to README hero position

## [0.9.0] - 2024-12-30

### Added

- **Anthropic API support** for session summarization as alternative to AWS Bedrock
  - Set `ANTHROPIC_API_KEY` to use direct API instead of Bedrock
- **Filter-only search** - search by date/project without requiring text query
  - `after:yesterday`, `before:2024-11-01`, `project:myapp` work standalone
- **SQLite busy timeout** - 5 second retry on database locks during concurrent access
- Message count and last working directory in CLI search output

### Fixed

- Date filter parsing: `2024-11-01` no longer misinterpreted as time "11:01"
- Hyphenated date filters now work: `3-days-ago`, `last-week`
- TUI search view overflow - header no longer disappears with many results
- Summary and project path truncation in TUI to prevent text spilling off screen
- Skip sessions with fewer than 5 messages during summarization

### Changed

- Improved summarization prompts for better problem-solution focus
- CLI search output format now matches TUI (relative time, message count, shorter paths)

## [0.2.6] - Previous release

See git history for earlier changes.
