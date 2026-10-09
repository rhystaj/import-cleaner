package languageprocessing

import (
	"bytes"
	"iter"
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

// DelimitedStringSequence divides a string into a sequence of substrings delimited by the specified runes.
// Each string in the sequence will contain the delimiting rune that terminated it.
func DelimitedStringSequence(s string, delimters []rune) iter.Seq[string] {
	var stringBuffer bytes.Buffer

	return func(yield func(string) bool) {
		for _, r := range s {
			if slices.Contains(delimters, r) {
				if !yield(stringBuffer.String() + string(r)) {
					return
				}
				stringBuffer.Reset()
				continue
			}

			stringBuffer.WriteRune(r)
		}

		if !yield(stringBuffer.String()) {
			return
		}
	}
}

// Source - https://stackoverflow.com/a/79034149
// Posted by eik, modified by community. See post 'Timeline' for change history
// Retrieved 2026-09-29, License - CC BY-SA 4.0

// Concat returns an iterator over the concatenation of the sequences.
func ConcatSeqs[V any](seqs ...iter.Seq[V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for e := range seq {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// SafelyDerefRuneToString returns the value of a rune pointer converted to a string if
// the pointer is not nil, or an empty string otherwise.
func SafelyDerefRuneToString(r *rune) string {
	if r == nil {
		return ""
	}

	return string(*r)
}
