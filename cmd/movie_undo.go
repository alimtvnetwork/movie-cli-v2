// movie_undo.go — movie movie undo
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mahin/mahin-cli-v2/db"
)

var movieUndoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the last move operation",
	Long:  `Reverts the most recent movie move operation.`,
	Run:   runMovieUndo,
}

func runMovieUndo(cmd *cobra.Command, args []string) {
	database, err := db.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Database error: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	lastMove, err := database.GetLastMove()
	if err != nil {
		fmt.Println("📭 No move operations to undo.")
		return
	}

	fmt.Println("⏪ Last move operation:")
	fmt.Println()
	fmt.Printf("  📁 %s\n", lastMove.ToPath)
	fmt.Printf("  → %s\n", lastMove.FromPath)
	fmt.Println()

	// Confirmation prompt
	fmt.Print("Undo this? [y/N]: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("❌ Undo cancelled.")
		return
	}

	// Check source exists
	if _, err := os.Stat(lastMove.ToPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "❌ File not found at: %s\n", lastMove.ToPath)
		fmt.Fprintln(os.Stderr, "   It may have been moved or deleted manually.")
		os.Exit(1)
	}

	// Move back
	if err := MoveFile(lastMove.ToPath, lastMove.FromPath); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Undo failed: %v\n", err)
		os.Exit(1)
	}

	// Mark as undone in DB
	if err := database.MarkMoveUndone(lastMove.ID); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not mark move as undone: %v\n", err)
	}

	// Update media path back
	if err := database.UpdateMediaPath(lastMove.MediaID, lastMove.FromPath); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Could not update media path: %v\n", err)
	}

	fmt.Println()
	fmt.Println("✅ Undo successful!")
	fmt.Printf("   %s\n", lastMove.ToPath)
	fmt.Printf("   → %s\n", lastMove.FromPath)
}
