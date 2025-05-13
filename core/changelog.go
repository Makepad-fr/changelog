package core

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Makepad-fr/semver/semver"
)

type Changelog struct {
	Unreleased Entry
	Releases   map[semver.Semver]Entry
}

func (c Changelog) HasUnreleased() bool {
	return len(c.Unreleased.Changes) > 0
}

type Entry struct {
	Date    string
	Changes map[Section][]string
}

type Section int

const (
	Added Section = iota
	Changed
	Deprecated
	Removed
	Fixed
	Security
)

func (s Section) String() string {
	switch s {
	case Added:
		return "Added"
	case Changed:
		return "Changed"
	case Deprecated:
		return "Deprecated"
	case Removed:
		return "Removed"
	case Fixed:
		return "Fixed"
	case Security:
		return "Security"
	default:
		return "Unknown"
	}
}

func ParseSection(s string) (Section, bool) {
	switch strings.ToLower(s) {
	case "added":
		return Added, true
	case "changed":
		return Changed, true
	case "deprecated":
		return Deprecated, true
	case "removed":
		return Removed, true
	case "fixed":
		return Fixed, true
	case "security":
		return Security, true
	default:
		return 0, false
	}
}

// Verify checks that version and date ordering is strictly descending.
// Returns an error if order is invalid.
func (c Changelog) Verify() error {
	type pair struct {
		version semver.Semver
		date    time.Time
	}

	var entries []pair
	for v, entry := range c.Releases {
		if entry.Date == "" {
			return fmt.Errorf("release date is empty for version %s", v)
		}

		t, err := time.Parse("2006-01-02", entry.Date)
		if err != nil {
			return fmt.Errorf("invalid date format for version %s: %v", v.String(), err)
		}
		entries = append(entries, pair{v, t})
	}

	// Sort descending by version
	sort.Slice(entries, func(i, j int) bool {
		return Compare(entries[i].version, entries[j].version) > 0
	})

	// Then: check that corresponding dates are also descending
	for i := 1; i < len(entries); i++ {
		vi := entries[i-1]
		vj := entries[i]
		if Compare(vi.version, vj.version) < 0 {
			return fmt.Errorf("version %s should come before %s", vj.version, vi.version)
		}
		if vi.date.Before(vj.date) {
			return fmt.Errorf("version %s has an older date (%s) than version %s (%s)",
				vi.version, vi.date.Format("2006-01-02"),
				vj.version, vj.date.Format("2006-01-02"))
		}
	}

	return nil
}

// GenerateMarkdown renders the Changelog in Keep a Changelog format.
func (c Changelog) GenerateMarkdown() string {
	var b strings.Builder

	b.WriteString("# Changelog\n\n")
	b.WriteString("All notable changes to this project will be documented in this file.\n")

	// Unreleased
	if c.HasUnreleased() {
		b.WriteString("## [Unreleased]")
		writeEntry(&b, c.Unreleased)
	}

	// Sort versions descending
	versions := make([]semver.Semver, 0, len(c.Releases))
	for v := range c.Releases {
		versions = append(versions, v)
	}
	sort.Slice(versions, func(i, j int) bool {
		return Compare(versions[i], versions[j]) > 0
	})

	for _, v := range versions {
		entry := c.Releases[v]
		line := fmt.Sprintf("\n## [%s]", v.String())
		if entry.Date != "" {
			line += fmt.Sprintf(" - %s", entry.Date)
		}
		b.WriteString(line + "\n")
		writeEntry(&b, entry)
	}

	return b.String()
}

func writeEntry(b *strings.Builder, entry Entry) {
	if len(entry.Changes) == 0 {
		b.WriteString("\n")
		return
	}

	for _, section := range []Section{
		Added, Changed, Deprecated, Removed, Fixed, Security,
	} {
		items, ok := entry.Changes[section]
		if !ok || len(items) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("\n### %s\n\n", section.String()))
		for _, item := range items {
			b.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
	}
}

func (cl *Changelog) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer f.Close()

	output := cl.GenerateMarkdown()
	if _, err := f.WriteString(output); err != nil {
		return fmt.Errorf("failed to write changelog content: %w", err)
	}

	return nil
}
