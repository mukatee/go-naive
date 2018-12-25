package chain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

//https://stackoverflow.com/questions/22811138/print-the-address-of-slice-in-golang

func TestTakeLongest(t *testing.T) {
	//create test chains, check if it changes to longest, ...
	GlobalChain = nil
	chain1 := createTestChain(10)
	GlobalChain = nil
	chain2 := createTestChain(16)
	assert.Equal(t, 11, len(chain1))
	assert.Equal(t, 17, len(chain2))
	GlobalChain = chain1
	assert.Equal(t, 11, len(GlobalChain))
	takeLongestChain(chain2)
	assert.Equal(t, 17, len(GlobalChain))
}

func TestTakeMostDifficult(t *testing.T) {
	//create test chains, check that it changes to the one with the highest difficulty
	GlobalChain = nil
	chain1 := createTestDiffChain(10, 1, 1, 1)
	GlobalChain = nil
	chain2 := createTestDiffChain(10, 1, 2, 3)
	diff1 := calculateChainDifficulty(chain1)
	diff2 := calculateChainDifficulty(chain2)
	GlobalChain = chain1
	println("diffs:", diff1, " ", diff2)
	takeMostDifficultChain(chain2)
	gDiff := calculateChainDifficulty(GlobalChain)
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
	for i := 0; i < size; i++ {
		data := fmt.Sprintf("Test%d", i)
		CreateBlock(nil, data, 0)
	}
	return GlobalChain
}

func createTestDiffChain(size int, diffs ...int) []Block {
	createTestChain(size)
	for i := 1; i <= size; i++ {
		GlobalChain[i].Difficulty = 10
		//previous hash is also used for block hash so have to re-set it before calculating hash
		GlobalChain[i].PreviousHash = GlobalChain[i-1].Hash
		hash := hash(&GlobalChain[i])
		GlobalChain[i].Hash = hash
	}
	println()
	for x := 0; x < len(diffs); x++ {
		idx := size + x + 1
		data := fmt.Sprintf("Test%d", x)
		CreateBlock(nil, data, 0)
		GlobalChain[idx].Difficulty = diffs[x]
		GlobalChain[idx].PreviousHash = GlobalChain[idx-1].Hash
		hash := hash(&GlobalChain[idx])
		GlobalChain[idx].Hash = hash
	}
	println()
	return GlobalChain
}
