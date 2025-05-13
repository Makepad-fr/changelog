package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Makepad-fr/semver/semver"
)

func Compare(a, b semver.Semver) int {
	aMajor := mustAtoi(a.Major)
	bMajor := mustAtoi(b.Major)
	if aMajor != bMajor {
		return compareInt(aMajor, bMajor)
	}

	aMinor := mustAtoi(a.Minor)
	bMinor := mustAtoi(b.Minor)
	if aMinor != bMinor {
		return compareInt(aMinor, bMinor)
	}

	aPatch := mustAtoi(a.Patch)
	bPatch := mustAtoi(b.Patch)
	if aPatch != bPatch {
		return compareInt(aPatch, bPatch)
	}

	return comparePreRelease(a.PreRelease, b.PreRelease)
}

func mustAtoi(s string) int {
	if s == "" {
		panic("semver field is empty string — invalid version component")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("invalid semver value %q: %v", s, err))
	}
	return n
}

func compareInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// Pre-release comparison: nil > pre-release
func comparePreRelease(a, b *string) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return 1 // release > pre-release
	}
	if b == nil {
		return -1 // pre-release < release
	}

	aParts := strings.Split(*a, ".")
	bParts := strings.Split(*b, ".")

	for i := 0; i < max(len(aParts), len(bParts)); i++ {
		var aPart, bPart string
		if i < len(aParts) {
			aPart = aParts[i]
		}
		if i < len(bParts) {
			bPart = bParts[i]
		}

		aNum, aErr := strconv.Atoi(aPart)
		bNum, bErr := strconv.Atoi(bPart)

		switch {
		case aErr == nil && bErr == nil:
			if aNum != bNum {
				return compareInt(aNum, bNum)
			}
		case aErr != nil && bErr != nil:
			if aPart != bPart {
				return strings.Compare(aPart, bPart)
			}
		case aErr == nil:
			return -1 // numeric < alphanumeric
		case bErr == nil:
			return 1
		}
	}

	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
