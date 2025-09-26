package services

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"retroblog/models"
	"strconv"
	"strings"
)

var (
	reNoteID          = regexp.MustCompile(`(?m)^Note ID:\s*(\d+)\s*$`)
	reTitle           = regexp.MustCompile(`(?m)^Title:\s*(.+)\s*$`)
	reFolder          = regexp.MustCompile(`(?m)^Folder:\s*(.+)\s*$`)
	reAttachCount     = regexp.MustCompile(`(?m)^Attachments:\s*(\d+)\s*file\(s\)\s*$`)
	reAttachmentStart = regexp.MustCompile(`(?i)^\s*\d+\.\s*Saved\s*as:\s*(.+?)\s*$`)
	reKeyVal          = regexp.MustCompile(`^\s*([A-Za-z ]+):\s*(.*?)\s*$`)
	reFirstInt        = regexp.MustCompile(`(\d+)`)
)

func ParseMetaFile(path string) (models.NoteMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return models.NoteMeta{}, err
	}
	return ParseMeta(string(data))
}

func ParseMeta(text string) (models.NoteMeta, error) {
	m := models.NoteMeta{Attachments: map[string]models.AttachmentMeta{}}
	t := normalizeNewlines(text)

	if id, err := findInt(reNoteID, t); err == nil {
		m.NoteID = id
	} else {
		return m, errors.New("meta missing Note ID")
	}

	if s, err := findString(reTitle, t); err == nil {
		m.Title = s
	} else {
		return m, errors.New("meta missing Title")
	}

	if s, err := findString(reFolder, t); err == nil {
		m.Folder = s
	} else {
		return m, errors.New("meta missing Folder")
	}

	if ms := reAttachCount.FindStringSubmatch(t); len(ms) == 2 {
		if v, err := strconv.Atoi(ms[1]); err == nil {
			m.AttachmentsDeclared = &v
		}
	}

	attachSec := extractAttachmentSection(t)
	if strings.TrimSpace(attachSec) == "" {
		return m, nil
	}

	lines := strings.Split(attachSec, "\n")
	var cur *models.AttachmentMeta

	flush := func() {
		if cur != nil {
			if cur.OriginalFilename != "" {
				m.Attachments[cur.OriginalFilename] = *cur
			} else if cur.SavedAs != "" {
				key := strings.TrimSuffix(cur.SavedAs, filepath.Ext(cur.SavedAs))
				m.Attachments[key] = *cur
			}
			cur = nil
		}
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		if ms := reAttachmentStart.FindStringSubmatch(line); len(ms) == 2 {
			flush()
			cur = &models.AttachmentMeta{SavedAs: strings.TrimSpace(ms[1])}
			continue
		}

		if cur != nil {
			if kv := reKeyVal.FindStringSubmatch(line); len(kv) == 3 {
				key := strings.ToLower(strings.TrimSpace(kv[1]))
				val := strings.TrimSpace(kv[2])
				switch key {
				case "original filename":
					cur.OriginalFilename = val
				case "source file":
					cur.SourceFile = val
				case "type":
					cur.Type = val
				case "position":
					if mm := reFirstInt.FindStringSubmatch(val); len(mm) == 2 {
						if iv, err := strconv.Atoi(mm[1]); err == nil {
							cur.Position = iv
						}
					}
				case "path":
					cur.Path = val
				case "id":
					cur.ID = val
				}
			}
		}
	}
	flush()

	return m, nil
}

func normalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

func findString(re *regexp.Regexp, s string) (string, error) {
	m := re.FindStringSubmatch(s)
	if len(m) != 2 {
		return "", errors.New("not found")
	}
	return strings.TrimSpace(m[1]), nil
}

func findInt(re *regexp.Regexp, s string) (int, error) {
	st, err := findString(re, s)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(st)
}

func extractAttachmentSection(s string) string {
	i := strings.Index(s, "Attachment Details:")
	if i == -1 {
		return ""
	}
	return s[i:]
}
