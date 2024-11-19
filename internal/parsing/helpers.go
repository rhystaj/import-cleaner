package parsing

import (
	"slices"
	"strings"
)

func IsOpenBracket(r rune) bool {
	return r == '(' || r == '{' || r == '[' || r == '<'
}

func IsCorrespondingCloseBracket(r rune, openBracket rune) bool {
	return (openBracket == '(' && r == ')') ||
		(openBracket == '{' && r == '}') ||
		(openBracket == '[' && r == ']') ||
		(openBracket == '<' && r == '>')
}

// SmallestCut preforms a cut at the earliest occurance of one of the specified runes.
func SmallestCut(s string, cutRunes []rune) (before string, after string, cutPerformed bool) {
	for _, r := range s {
		if slices.Contains(cutRunes, r) {
			return strings.Cut(s, string(r))
		}
	}

	return s, "", false
}
