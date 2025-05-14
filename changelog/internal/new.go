package internal

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Makepad-fr/changelog/core"
	"github.com/Makepad-fr/changelog/parser"
	"github.com/Makepad-fr/semver/semver"
)

func NewCommand() {
	newFlags := flag.NewFlagSet("new", flag.ExitOnError)
	var added, changed, removed, fixed, security StringSliceFlag
	noPrompt := newFlags.Bool("no-prompt", false, "Disable all interactive prompts")

	newFlags.Var(&added, "added", "Specify multiple --added flags for multiple items")
	newFlags.Var(&changed, "changed", "Specify multiple --changed flags for multiple items")
	newFlags.Var(&removed, "removed", "Specify multiple --removed flags for multiple items")
	newFlags.Var(&fixed, "fixed", "Specify multiple --fixed flags for multiple items")
	newFlags.Var(&security, "security", "Specify multiple --security flags for multiple items")

	file := newFlags.String("file", "CHANGELOG.md", "Path to the changelog file")
	version := newFlags.String("version", "", "Version number (e.g. 1.2.3)")

	newFlags.Parse(os.Args[2:])

	if *version == "" {
		if *noPrompt {
			log.Fatalln("❌ --version is required when --no-prompt is enabled.")
		}
		fmt.Print("Enter version (e.g., 1.2.3): ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			*version = strings.TrimSpace(scanner.Text())
		}
		if *version == "" {
			log.Fatalln("❌ Version is required.")
		}
	}

	v, err := semver.Parse(*version)
	if err != nil {
		log.Fatalf("❌ Invalid version %q: %v", *version, err)
	}

	if *noPrompt {
		total := len(added) + len(changed) + len(removed) + len(fixed) + len(security)
		if *version == "" {
			log.Fatalln("❌ --version is required when --no-prompt is enabled.")
		}
		if total == 0 {
			log.Fatalln("❌ At least one change section (e.g. --added) must be defined when --no-prompt is enabled.")
		}
	}

	var cl *core.Changelog
	var existingEntry core.Entry
	versionExists := false

	if _, err := os.Stat(*file); err == nil {
		loaded, err := parser.LoadChangelogFromFile(*file)
		if err != nil {
			log.Fatalf("Failed to load changelog: %v", err)
		}
		cl = loaded
		if entry, ok := cl.Releases[v]; ok {
			versionExists = true
			existingEntry = entry
		}
	} else if os.IsNotExist(err) {
		cl = &core.Changelog{
			Unreleased: core.Entry{Changes: make(map[core.Section][]string)},
			Releases:   map[semver.Semver]core.Entry{},
		}
	} else {
		log.Fatalf("Failed to check changelog file: %v", err)
	}

	if versionExists {
		if *noPrompt {
			log.Fatalf("❌ Version %s already exists. Cannot continue in --no-prompt mode.", v.String())
		}
		fmt.Printf("⚠️ Version %s already exists. Choose action:\n", v.String())
		fmt.Println("1) Override")
		fmt.Println("2) Update (merge)")
		fmt.Println("3) Exit")
		fmt.Print("Enter choice [1/2/3]: ")

		scanner := bufio.NewScanner(os.Stdin)
		var choice string
		if scanner.Scan() {
			choice = strings.TrimSpace(scanner.Text())
		}

		switch choice {
		case "1":
			// continue to override
		case "2":
			mergeSection(existingEntry.Changes, core.Added, added, *noPrompt)
			mergeSection(existingEntry.Changes, core.Changed, changed, *noPrompt)
			mergeSection(existingEntry.Changes, core.Removed, removed, *noPrompt)
			mergeSection(existingEntry.Changes, core.Fixed, fixed, *noPrompt)
			mergeSection(existingEntry.Changes, core.Security, security, *noPrompt)

			existingEntry.Date = time.Now().Format("2006-01-02")
			cl.Releases[v] = existingEntry
			cl.WriteToFile(*file)
			return
		case "3", "":
			fmt.Println("❌ Aborted.")
			os.Exit(0)
		default:
			log.Fatalln("❌ Invalid choice.")
		}
	}

	changes := map[core.Section][]string{}
	addSection(changes, core.Added, added, *noPrompt)
	addSection(changes, core.Changed, changed, *noPrompt)
	addSection(changes, core.Removed, removed, *noPrompt)
	addSection(changes, core.Fixed, fixed, *noPrompt)
	addSection(changes, core.Security, security, *noPrompt)

	cl.Releases[v] = core.Entry{
		Date:    time.Now().Format("2006-01-02"),
		Changes: changes,
	}

	cl.WriteToFile(*file)
}

func addSection(changes map[core.Section][]string, section core.Section, value []string, noPrompt bool) {
	if len(value) == 0 {
		if noPrompt {
			changes[section] = []string{} // consider empty silently
			return
		}
		value = promptItems(section.String())
	}
	changes[section] = value
}

func mergeSection(existing map[core.Section][]string, section core.Section, value []string, noPrompt bool) {
	if len(value) == 0 {
		if noPrompt {
			return // skip if not provided
		}
		value = promptItems(section.String())
	}
	existing[section] = append(existing[section], value...)
}

func promptItems(title string) []string {
	fmt.Printf("Enter %s items (press Enter twice to finish):\n", title)
	scanner := bufio.NewScanner(os.Stdin)
	var items []string
	for {
		fmt.Printf("- ")
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			break
		}
		items = append(items, text)
	}
	return items
}
