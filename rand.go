package profaneword

import (
	"crypto/rand"
	"math"
	"math/big"
)

type RandomDevice interface {
	//Rand is a function that returns a random number between 0 and 1
	Rand() *big.Rat
	RandMax(max int64) int64
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

func (c CryptoRand) RandMax(max int64) int64 {
	return getRand(big.NewInt(max)).Int64()
}
