package profanities

type Word uint8

const (
	//NONE is the default: there is no word at this RadixWordNode
	NONE  Word = 0
	START Word = 1 << iota
	FILLER
	END
	EXCL
	//SPLIT reserved for future sentence-construction
	SPLIT
	//MISSPELL is inherited in the tree; covers slang and miss-spelling
	MISSPELL
	//POSITIVE is also inherited, and cover words that are not negatively laden
	POSITIVE
	//DEFAULT covers most normal-kinda words
	DEFAULT = START | FILLER
	EXCLS   = START | EXCL
)