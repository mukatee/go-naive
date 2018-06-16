package chain

import "strings"

var DIFFICULTY_ADJUSTMENT_INTERVAL = 10 //number of blocks between to aim to adjust difficulty
var BLOCK_GENERATION_INTERVAL = 10 //target seconds to generate a block

func verifyHashVsDifficulty(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

func getDifficulty() int {
	prevBlock := globalChain[len(globalChain)-1]
	if (prevBlock.Index % DIFFICULTY_ADJUSTMENT_INTERVAL == 0 && prevBlock.Index != 0) {
		return getAdjustedDifficulty(prevBlock, globalChain);
	} else {
		return prevBlock.Difficulty;
	}
	return 0
}

func getAdjustedDifficulty(prevBlock Block, chain []Block) int {
	prevAdjustmentBlock := chain[len(chain) - DIFFICULTY_ADJUSTMENT_INTERVAL];
	timeExpected := BLOCK_GENERATION_INTERVAL * DIFFICULTY_ADJUSTMENT_INTERVAL;
	timeDelta := prevBlock.Timestamp.Sub(prevAdjustmentBlock.Timestamp);
	seconds := int(timeDelta.Seconds())
	if (seconds < timeExpected / 2) {
		return prevAdjustmentBlock.Difficulty + 1;
	} else if (seconds > timeExpected * 2) {
		return prevAdjustmentBlock.Difficulty - 1;
	} else {
		return prevAdjustmentBlock.Difficulty;
	}
	return 0
}

func validateTimestamp(newBlock Block, prevBlock Block) bool {
	diff := newBlock.Timestamp.Sub(prevBlock.Timestamp)
	pastOk := diff > -60
	futureOk := diff < 60
	return pastOk && futureOk
}

