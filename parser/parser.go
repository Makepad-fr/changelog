package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/Makepad-fr/changelog/core"
	"github.com/Makepad-fr/semver/semver"
)

// Parse parses a changelog file and returns a structured Changelog.
func Parse(r io.Reader) (*core.Changelog, error) {
	reader := bufio.NewReader(r)
	cl := &core.Changelog{Releases: make(map[semver.Semver]core.Entry)}
	currentEntry := core.Entry{Changes: make(map[core.Section][]string)}
	currentVersion := ""
	currentSection := core.Section(0)
	inUnreleased := false

	reRelease := regexp.MustCompile(`^## \[(.*?)\](?: - (\d{4}-\d{2}-\d{2}))?`)
	reSection := regexp.MustCompile(`^### (.+)`)
	reBullet := regexp.MustCompile(`^- (.+)`)
	var buffer bytes.Buffer
	var line string
	for {
		chunk, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				if buffer.Len() > 0 {
					line = buffer.String()
					buffer.Reset()
					// Process last line before breaking
				} else {
					break // ← this finally exits the loop
				}
			} else {
				return nil, fmt.Errorf("read error: %w", err)
			}
		} else {
			buffer.Write(chunk)
			if isPrefix {
				continue // line is too long, keep reading
			}
			line = buffer.String()
			buffer.Reset()
		}
		line = strings.ReplaceAll(line, "\u00A0", " ")
		matches := reRelease.FindStringSubmatch(line)
		if len(matches) > 0 {
			if currentVersion != "" {
				if inUnreleased {
					cl.Unreleased = currentEntry
				} else {
					v, err := semver.Parse(currentVersion)
					if err == nil {
						cl.Releases[v] = currentEntry
					}
				}
			}

			currentEntry = core.Entry{Changes: make(map[core.Section][]string)}
			currentVersion = matches[1]
			inUnreleased = strings.ToLower(currentVersion) == "unreleased"

			if len(matches) > 2 {
				currentEntry.Date = matches[2]
			}
			continue
		}
		matches = reSection.FindStringSubmatch(line)
		if len(matches) > 0 {
			if sec, ok := core.ParseSection(matches[1]); ok {
				currentSection = sec
				if _, exists := currentEntry.Changes[currentSection]; !exists {
					currentEntry.Changes[currentSection] = []string{}
				}
			}
			continue
		}
		matches = reBullet.FindStringSubmatch(line)
		if len(matches) > 0 {
			if _, exists := currentEntry.Changes[currentSection]; exists {
				currentEntry.Changes[currentSection] = append(currentEntry.Changes[currentSection], matches[1])
			}
		}
	}
	if currentVersion != "" {
		if inUnreleased {
			cl.Unreleased = currentEntry
		} else {
			v, err := semver.Parse(currentVersion)
			if err == nil {
				cl.Releases[v] = currentEntry
			}
		}
	}

	return cl, nil
}

func LoadChangelogFromFile(path string) (*core.Changelog, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}
