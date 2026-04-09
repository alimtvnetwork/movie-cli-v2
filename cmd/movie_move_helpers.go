// movie_move_helpers.go — helper functions for movie move command
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mahin/mahin-cli-v2/cleaner"
	"github.com/mahin/mahin-cli-v2/db"
)

// promptSourceDirectory shows configured directories and a custom option.
func promptSourceDirectory(scanner *bufio.Scanner, database *db.DB, home string) string {
	scanDir, err := database.GetConfig("scan_dir")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error: %v\n", err)
	}
	scanDir = expandHome(scanDir, home)

	fmt.Println("📂 Choose a directory to browse:")
	fmt.Println()
	fmt.Printf("  1. 📥 Downloads   (%s)\n", expandHome("~/Downloads", home))
	if scanDir != "" && scanDir != expandHome("~/Downloads", home) {
		fmt.Printf("  2. 🔍 Scan Dir    (%s)\n", scanDir)
		fmt.Println("  3. 📁 Custom path")
		fmt.Println()
		fmt.Print("  Choose [1-3]: ")
	} else {
		fmt.Println("  2. 📁 Custom path")
		fmt.Println()
		fmt.Print("  Choose [1-2]: ")
	}

	if !scanner.Scan() {
		return ""
	}

	input := strings.TrimSpace(scanner.Text())
	hasScanDir := scanDir != "" && scanDir != expandHome("~/Downloads", home)

	switch input {
	case "1":
		return expandHome("~/Downloads", home)
	case "2":
		if hasScanDir {
			return scanDir
		}
		return promptCustomPath(scanner, home)
	case "3":
		if hasScanDir {
			return promptCustomPath(scanner, home)
		}
		fmt.Println("❌ Invalid choice")
		return ""
	default:
		fmt.Println("❌ Invalid choice")
		return ""
	}
}

func promptCustomPath(scanner *bufio.Scanner, home string) string {
	fmt.Print("  Enter path: ")
	if !scanner.Scan() {
		return ""
	}
	return expandHome(strings.TrimSpace(scanner.Text()), home)
}

// promptDestination shows destination options.
func promptDestination(scanner *bufio.Scanner, database *db.DB, home string) string {
	moviesDir, err := database.GetConfig("movies_dir")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error: %v\n", err)
	}
	tvDir, err := database.GetConfig("tv_dir")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error: %v\n", err)
	}
	archiveDir, err := database.GetConfig("archive_dir")
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Config read error: %v\n", err)
	}

	moviesDir = expandHome(moviesDir, home)
	tvDir = expandHome(tvDir, home)
	archiveDir = expandHome(archiveDir, home)

	fmt.Println()
	fmt.Println("  Destination:")
	fmt.Printf("  1. 🎬 Movies     (%s)\n", moviesDir)
	fmt.Printf("  2. 📺 TV Shows   (%s)\n", tvDir)
	fmt.Printf("  3. 📦 Archive    (%s)\n", archiveDir)
	fmt.Println("  4. 📁 Custom path")
	fmt.Println()
	fmt.Print("  Choose [1-4]: ")

	if !scanner.Scan() {
		return ""
	}

	switch strings.TrimSpace(scanner.Text()) {
	case "1":
		return moviesDir
	case "2":
		return tvDir
	case "3":
		return archiveDir
	case "4":
		return promptCustomPath(scanner, home)
	default:
		fmt.Println("❌ Invalid choice")
		return ""
	}
}

// listVideoFiles returns video files in a directory sorted by name.
func listVideoFiles(dir string) []os.FileInfo {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var files []os.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if cleaner.IsVideoFile(entry.Name()) {
			info, err := entry.Info()
			if err == nil {
				files = append(files, info)
			}
		}
	}
	return files
}

// humanSize formats bytes into human-readable size.
func humanSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func expandHome(path, home string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	if path == "~" {
		return home
	}
	return path
}

func saveHistoryLog(basePath, title string, year int, from, to string) {
	slug := cleaner.ToSlug(title)
	if year > 0 {
		slug += "-" + strconv.Itoa(year)
	}
	histDir := filepath.Join(basePath, "json", "history", slug)
	if err := os.MkdirAll(histDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Cannot create history dir: %v\n", err)
		return
	}

	logFile := filepath.Join(histDir, "move-log.json")

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Cannot open history log: %v\n", err)
		return
	}
	defer f.Close()

	entry := fmt.Sprintf(`{"from":"%s","to":"%s","timestamp":"%s"}`+"\n",
		from, to, time.Now().Format(time.RFC3339))
	if _, err := f.WriteString(entry); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Cannot write history log: %v\n", err)
	}
}

// MoveFile moves a file from src to dst. It first attempts os.Rename (atomic,
// same-filesystem). If that fails with EXDEV (cross-device link), it falls back
// to copy-then-remove. The source is only deleted after a successful copy+close.
//
// SHARED: used by move, rename, undo
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	if !isCrossDeviceError(err) {
		return fmt.Errorf("rename failed: %w", err)
	}

	// Cross-device fallback: copy + remove
	return crossDeviceMove(src, dst)
}

// isCrossDeviceError checks whether the error is an EXDEV (cross-device link)
// error, which occurs when os.Rename is called across different filesystems
// (e.g., USB drives, network mounts, different partitions).
func isCrossDeviceError(err error) bool {
	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		if errno, ok := linkErr.Err.(syscall.Errno); ok {
			return errno == syscall.EXDEV
		}
	}
	return false
}

// crossDeviceMove copies the file from src to dst, preserves the original file
// permissions, and removes the source only after the destination is fully
// written and closed. If anything fails, the source file is left untouched.
func crossDeviceMove(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("stat source: %w", err)
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("create destination: %w", err)
	}

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		dstFile.Close()
		os.Remove(dst) // clean up partial file
		return fmt.Errorf("copy failed: %w", err)
	}

	if err = dstFile.Close(); err != nil {
		os.Remove(dst)
		return fmt.Errorf("close destination: %w", err)
	}

	// Source file is only removed after successful copy + close
	if err = os.Remove(src); err != nil {
		return fmt.Errorf("remove source (file copied successfully to %s): %w", dst, err)
	}

	return nil
}
