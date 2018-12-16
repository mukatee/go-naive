package chain

import (
	"log"
	"strings"
)

var DIFFICULTY_ADJUSTMENT_INTERVAL = 10 //number of blocks between to aim to adjust difficulty
var BLOCK_GENERATION_INTERVAL = 10      //target seconds to generate a block

//verifyHashVsDifficulty checks if given hash string starts with "difficulty" number of zeroes
//TODO: check against difficulty number so hash < difficulty for much more granularity
func verifyHashVsDifficulty(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

//getDifficulty calculates the current difficulty based on timestamps
func getDifficulty() int {
	prevBlock := GlobalChain[len(GlobalChain)-1]
	if prevBlock.Index%DIFFICULTY_ADJUSTMENT_INTERVAL == 0 && prevBlock.Index != 0 {
		return getAdjustedDifficulty(prevBlock, GlobalChain)
	} else {
		return prevBlock.Difficulty
	}
	return 0
}

//getAdjustedDifficulty increases or decreases difficulty based on time diff in finding new block
func getAdjustedDifficulty(prevBlock Block, chain []Block) int {
	prevAdjustmentBlock := chain[len(chain)-DIFFICULTY_ADJUSTMENT_INTERVAL]
	timeExpected := BLOCK_GENERATION_INTERVAL * DIFFICULTY_ADJUSTMENT_INTERVAL
	timeDelta := prevBlock.Timestamp.Sub(prevAdjustmentBlock.Timestamp)
	seconds := int(timeDelta.Seconds())
	if seconds < timeExpected/2 {
		log.Print("Difficulty less than half of expected: ", seconds, " vs. ", timeExpected, " -> reducing.")
		//less than half expected time, increase difficulty by one
		return prevAdjustmentBlock.Difficulty + 1
	} else if seconds > timeExpected*2 {
		log.Print("Difficulty twice the expected: ", seconds, " vs. ", timeExpected, " -> increasing.")
		//more than twice expected time, decrease difficulty by one
		return prevAdjustmentBlock.Difficulty - 1
	} else {
		return prevAdjustmentBlock.Difficulty
	}
	return 0
}

//validateTimestamp checks that block timestamp is withing the allowed time interval from previous block
func validateTimestamp(newBlock Block, prevBlock Block) bool {
	diff := newBlock.Timestamp.Sub(prevBlock.Timestamp)
	pastOk := diff > -60
	futureOk := diff < 60
	return pastOk && futureOk
}
