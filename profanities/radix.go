package profanities

type radixWordNode struct {
	val      string
	word     Word
	branches []*radixWordNode
}

func (n *radixWordNode) getWordsOf(words []Word, dissallowedWord Word) map[Word][]string {
	return n.getWords(``, words, dissallowedWord)
}

func (n *radixWordNode) getWords(base string, words []Word, dissallowedWord Word) map[Word][]string {
	if n.word&dissallowedWord != 0 {
		return nil
	}
	mp := make(map[Word][]string, len(words))
	text := base + n.val
	for _, word := range words {
		if n.word&word != 0 {
			mp[word] = append(mp[word], text)
		}
	}
	for _, branch := range n.branches {
		for w, strs := range branch.getWords(text, words, dissallowedWord) {
			mp[w] = append(mp[w], strs...)
		}
	}
	return mp
}

//GetOfSingle returns the Word's of a single Word type, and subtracts all words of the dissallowed type, Word
func (n *radixWordNode) GetOfSingle(word, dissallowedWord Word) []string {
	return n.getWordsOf([]Word{word}, dissallowedWord)[word]
}
