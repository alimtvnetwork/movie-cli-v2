# Changelog

All notable changes to this project will be documented in this file.

## v0.1.0

### Added
- Core CLI with Cobra: `hello`, `version`, `self-update` commands
- Movie management: `scan`, `ls`, `search`, `info`, `suggest`, `move`, `rename`, `undo`, `play`, `stats`, `config`, `tag`
- SQLite database with WAL mode, 5 tables, 7 indexes
- TMDb API client (search, details, credits, recommendations, trending, posters)
- Filename cleaner (junk removal, year extraction, TV detection)
- Cross-drive move fallback (copy+delete when os.Rename fails)
- PowerShell build & deploy pipeline (`run.ps1`)
- GitHub Actions release pipeline with cross-compiled binaries
- Cross-platform install scripts (install.ps1, install.sh)
- Comprehensive `--help` on all commands
- Full project specification in `spec/`
