package profanities

import (
	"fmt"
	"testing"
)

type aggregate struct {
	count          int
	averageWordLen float32
}

func (a aggregate) combine(o aggregate) aggregate {
	return aggregate{
		count:          a.count + o.count,
		averageWordLen: (a.averageWordLen*float32(a.count) + o.averageWordLen*float32(o.count)) / float32(a.count+o.count),
	}
}

func (a aggregate) multiply(o aggregate) aggregate {
	tmp := a.combine(o)
	tmp.count = a.count * o.count
	return tmp
}

func (a *aggregate) setAverage(total int) {
	a.averageWordLen = float32(total) / float32(a.count)
}

type aggregator map[Word]aggregate

func (a aggregator) getAggregate(word Word) aggregate {
	if agg, found := a[word]; found {
		return agg
	}
	a[word] = radixDatabaseAggr(word)
	return a[word]
}

func radixDatabaseAggr(wordToAggregate Word) aggregate {
	countMap := make(map[string]struct{}, len(wordData))
	charsTotal := 0
	for _, r := range wordData {
		words := r.GetOfSingle(wordToAggregate, WEIRD)
		for _, wrd := range words {
			if _, found := countMap[wrd]; !found {
				countMap[wrd] = struct{}{}
				charsTotal += len(wrd)
			}
		}
	}
	a := aggregate{count: len(countMap)}
	a.setAverage(charsTotal)
	return a
}

func TestReport(_ *testing.T) {
	aggr := make(aggregator)
	lastWord := aggregate{}
	otherWord := aggregate{}
	for _, s := range sentences {
		sentA := aggr.getAggregate(s.sentnc.word)
		otherWord = otherWord.combine(sentA)
		if s.sentPos == 0 {
			lastWord = lastWord.combine(sentA)
		}
	}

	for i := range [3]struct{}{} {
		agg := lastWord
		for j := 0; j < i; j++ {
			agg = agg.multiply(otherWord)
		}
		fmt.Printf("for %d word there will be %d combinations and on average %.2f letters in each word\n", i+1, agg.count, agg.averageWordLen)
	}
}
