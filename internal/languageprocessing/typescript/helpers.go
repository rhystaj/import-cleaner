package typescript

import (
	"bytes"
	"importcleaner/internal/datastructures"
	"iter"
)

func stripRelativePrefixFromModulePath(modulePath []string) []string {
	i := 0
	for i < len(modulePath) && (modulePath[i] == "." || modulePath[i] == "..") {
		i++
	}
	return modulePath[i:]
}

// Returns a sequence of fragments of the input text split by Typescript statement terminators in the order they appear.
// A Typescript terminator is a semicolon followed by a new line, or a semicolon or newline on their own.
func fragmentByTypescriptStatementTerminators(text string) iter.Seq[string] {
	var stringBuffer bytes.Buffer
	lookahead, _ := datastructures.CreateNewLookahead[rune](2)

	return func(yield func(string) bool) {
		ignoreNextRune := false
		for lah := range lookahead.ProcessSequence(stringAsRuneSequence(text)) {
			if ignoreNextRune {
				ignoreNextRune = false
				continue
			}

			if lah.PeekStart() == nil {
				// In theory this shouldn't occur due to how ProcessSequence is implemented - check is defensive
				continue
			}

			currentRune := *lah.PeekStart()
			stringBuffer.WriteRune(currentRune)

			if currentRune != '\n' && currentRune != ';' {
				continue
			}

			if currentRune == ';' && lah.PeekEnd() != nil && *lah.PeekEnd() == '\n' {
				ignoreNextRune = true
				stringBuffer.WriteRune('\n')
			}

			if !yield(stringBuffer.String()) {
				return
			}
			stringBuffer.Reset()
		}

		if stringBuffer.Len() > 0 {
			if !yield(stringBuffer.String()) {
				return
			}
		}
	}

}
