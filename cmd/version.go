// version.go — implements the `mahin version` command.
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/mahin/mahin-cli-v1/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current version, commit, and build date",
	Long: `Display the full version information for the mahin binary.

Shows the semantic version, git commit hash, build date, Go version,
and OS/architecture. Useful for debugging and reporting issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mahin %s\n", version.Full())
		fmt.Printf("  Go:   %s\n", runtime.Version())
		fmt.Printf("  OS:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}
