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

type CryptoRand struct {}

func (c CryptoRand) Rand() *big.Rat {
	max := big.NewInt(math.MaxInt64)
	out,_ := rand.Int(rand.Reader, max)
	return big.NewRat(out.Int64(), max.Int64())
}

func (c CryptoRand) RandMax(max int64) int64 {
	bigMax := big.NewInt(max)
	out,_ := rand.Int(rand.Reader, bigMax)
	return out.Int64()
}