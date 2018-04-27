package chain

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

//https://stackoverflow.com/questions/22811138/print-the-address-of-slice-in-golang

func TestTakeLongest(t *testing.T) {
	//create test chains, check if it changes to longest, ...
	globalChain = nil
	chain1 := createTestChain(10)
	globalChain = nil
	chain2 := createTestChain(16)
	assert.Equal(t, 11, len(chain1))
	assert.Equal(t, 17, len(chain2))
	globalChain = chain1
	assert.Equal(t, 11,len(globalChain))
	takeLongestChain(chain2)
	assert.Equal(t, 17, len(globalChain))
}

func TestTakeMostDifficult(t *testing.T) {
	globalChain = nil
	chain1 := createTestDiffChain(10, 1, 1, 1)
	globalChain = nil
	chain2 := createTestDiffChain(10, 1, 2, 3)
	diff1 := calculateChainDifficulty(chain1)
	diff2 := calculateChainDifficulty(chain2)
	globalChain = chain1
	println("diffs:", diff1, " ", diff2)
	takeMostDifficultChain(chain2)
	gDiff := calculateChainDifficulty(globalChain)
	println("diffs:", diff1, " ", diff2, " ", gDiff)
	if diff1 > diff2 {
		assert.Equal(t, diff1, gDiff)
	} else {
		assert.Equal(t, diff2, gDiff)
	}
//	diff1Str := strconv.FormatFloat(diff1, 'f', 2, 64)
//	println("diff1: ", diff1Str)
//	diff2Str := strconv.FormatFloat(diff2, 'f', 2, 64)
//	println("diff2: ", diff2Str)
}

//create a test chain of given length
func createTestChain(size int) []Block {
	createGenesisBlock(true)
	for i := 0 ; i < size ; i++ {
		data := fmt.Sprintf("Test%d", i)
		createBlock(nil, data, 0)
	}
	return globalChain
}

func createTestDiffChain(size int, diffs ...int) []Block {
	createTestChain(size)
	for i := 1 ; i <= size ; i++ {
		globalChain[i].Difficulty = 10
		//previous hash is also used for block hash so have to re-set it before calculating hash
		globalChain[i].PreviousHash = globalChain[i-1].Hash
		hash := hash(&globalChain[i])
		globalChain[i].Hash = hash
	}
	println()
	for x := 0 ; x < len(diffs) ; x++ {
		idx := size + x + 1
		data := fmt.Sprintf("Test%d", x)
		createBlock(nil, data, 0)
		globalChain[idx].Difficulty = diffs[x]
		globalChain[idx].PreviousHash = globalChain[idx-1].Hash
		hash := hash(&globalChain[idx])
		globalChain[idx].Hash = hash
	}
	println()
	return globalChain
}
















