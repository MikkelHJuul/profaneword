package profaneword

import (
	"crypto/rand"
	"math"
	"math/big"
)

type RandomDevice interface {
	// Rand is a function that returns a random number between 0 and 1
	Rand() *big.Rat
	RandMax(max int) int
}

func getRand(max *big.Int) *big.Int {
	i, _ := rand.Int(rand.Reader, max)
	return i
}

type CryptoRand struct{}

func (c CryptoRand) Rand() *big.Rat {
	max := big.NewInt(math.MaxInt64)
	return big.NewRat(getRand(max).Int64(), max.Int64())
}

func (c CryptoRand) RandMax(max int) int {
	return int(getRand(big.NewInt(int64(max))).Int64())
}

type thresholdRandom struct {
	Rand      RandomDevice
	Threshold *big.Rat
}

func newRandomFormatter(device RandomDevice, threshold *big.Rat) thresholdRandom {
	if threshold == nil {
		threshold = big.NewRat(1, 2)
	}
	if device == nil {
		device = CryptoRand{}
	}
	return thresholdRandom{
		Rand:      device,
		Threshold: threshold,
	}
}

func newFiftyFifty() thresholdRandom {
	return newRandomFormatter(nil, nil)
}
