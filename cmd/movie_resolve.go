// movie_resolve.go — shared media resolver
//
// ── Shared helper exported from this file ───────────────────────────
//
//   resolveMediaByQuery(db, query) (*Media, error)
//       Resolves a media item from the local DB by numeric ID or fuzzy
//       title match (exact → prefix → first result).
//
// Consumers: movie_info.go, movie_play.go, movie_ls.go (detail view)
//
// All commands that accept an <id-or-title> argument should use this
// helper to keep resolution logic consistent.  Do NOT duplicate the
// ID-parse → exact-match → prefix-match → fallback chain elsewhere.
// ────────────────────────────────────────────────────────────────────
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mahin/mahin-cli-v2/db"
)

// resolveMediaByQuery resolves a media item by numeric ID or fuzzy title query.
func resolveMediaByQuery(database *db.DB, query string) (*db.Media, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("empty media identifier")
	}

	if id, err := strconv.ParseInt(query, 10, 64); err == nil {
		m, err := database.GetMediaByID(id)
		if err != nil {
			return nil, fmt.Errorf("media not found for ID %d", id)
		}
		return m, nil
	}

	results, err := database.SearchMedia(query)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("media not found for %q", query)
	}

	for _, m := range results {
		if strings.EqualFold(m.CleanTitle, query) || strings.EqualFold(m.Title, query) {
			picked := m
			return &picked, nil
		}
	}

	queryLower := strings.ToLower(query)
	for _, m := range results {
		if strings.HasPrefix(strings.ToLower(m.CleanTitle), queryLower) ||
			strings.HasPrefix(strings.ToLower(m.Title), queryLower) {
			picked := m
			return &picked, nil
		}
	}

	picked := results[0]
	return &picked, nil
}
