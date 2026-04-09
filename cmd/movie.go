// movie.go — parent command: movie movie
package cmd

import (
	"github.com/spf13/cobra"
)

var movieCmd = &cobra.Command{
	Use:   "movie",
	Short: "Movie & TV show manager",
	Long: `Manage your personal movie and TV show library.

Scanning & Metadata:
  scan <folder>           Scan a folder for video files and fetch TMDb metadata
  search <query>          Search TMDb for movies/TV shows and save to library
  info <id|title>         Show detailed info (local DB → TMDb fallback)

Browsing:
  ls                      Paginated list of your library
  stats                   Library statistics (counts, genres, ratings)
  suggest                 Get recommendations and trending titles

File Management:
  move                    Interactively browse and move video files
  rename                  Batch-rename messy filenames using clean format
  undo                    Revert the last move or rename operation
  play <id>               Open a video with your system's default player

Organization:
  tag add <id> <tag>      Add a tag to a media item
  tag remove <id> <tag>   Remove a tag from a media item
  tag list [id]           List tags for an item or all tags

Configuration:
  config get <key>        View a configuration value
  config set <key> <val>  Set a configuration value (e.g., tmdb-key)

Examples:
  movie movie scan ~/Movies
  movie movie ls --page 2
  movie movie search "The Matrix"
  movie movie info 42
  movie movie suggest --genre action
  movie movie tag add 1 favorite`,
}

func init() {
	movieCmd.AddCommand(
		movieScanCmd,
		movieLsCmd,
		movieSearchCmd,
		movieSuggestCmd,
		movieMoveCmd,
		movieUndoCmd,
		movieInfoCmd,
		moviePlayCmd,
		movieStatsCmd,
		movieRenameCmd,
		movieConfigCmd,
	)
}
