package main

import (
	"regexp"
	"strings"
)

// Helper function to count protected posts
func CountProtectedPosts(posts []post) int {
	count := 0
	for _, p := range posts {
		if p.IsProtected {
			count++
		}
	}
	return count
}

// ExtractDescription extracts description from content
func ExtractDescription(body string) string {
	desc := body
	desc = regexp.MustCompile(`(?m)^#+\s+.*$`).ReplaceAllString(desc, "")
	desc = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`).ReplaceAllString(desc, "$1")
	desc = regexp.MustCompile(`!\[([^\]]*)\]\([^\)]+\)`).ReplaceAllString(desc, "")
	desc = regexp.MustCompile(`\*+([^\*]+)\*+`).ReplaceAllString(desc, "$1")
	desc = regexp.MustCompile(`_+([^_]+)_+`).ReplaceAllString(desc, "$1")
	desc = regexp.MustCompile("(?s)```[^`]*```").ReplaceAllString(desc, "")
	desc = regexp.MustCompile("`([^`]+)`").ReplaceAllString(desc, "$1")
	desc = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(desc, "")
	desc = regexp.MustCompile(`\[Attachment:[^\]]+\]`).ReplaceAllString(desc, "")
	desc = strings.TrimSpace(desc)
	desc = regexp.MustCompile(`\s+`).ReplaceAllString(desc, " ")

	if len(desc) > 160 {
		if idx := strings.LastIndex(desc[:157], " "); idx > 100 {
			desc = desc[:idx] + "..."
		} else {
			desc = desc[:157] + "..."
		}
	}

	return desc
}
