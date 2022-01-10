package profanities

// Word is a bitmask for marking Word type in this library, fx, placement in a sentence
type Word uint8

const (
	// START at the start of a word
	START Word = 1 << iota
	// FILLER - just somewhere
	FILLER
	// END of the sentence
	END
	// EXCL - as an exclamation, like: "DAMN!"
	EXCL
	// SPLIT reserved for future sentence-construction (splitting using words)
	SPLIT
	// MISSPELL is inherited in the tree; covers slang and miss-spelling
	MISSPELL
	// POSITIVE is also inherited, and cover words that are not negatively laden
	POSITIVE
	// DEFAULT covers most normal-kinda words
	DEFAULT = START | FILLER
	// EXCLS just a concatenation, because they appear together often
	EXCLS = START | EXCL
	// NONE is the default: there is no word at this radixWordNode
	NONE Word = 0
)
