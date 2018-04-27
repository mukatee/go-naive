package tests

import (
	"testing"
	"time"
	"github.com/mukatee/go-naive/chain"
	"github.com/stretchr/testify/assert"
)

//https://stackoverflow.com/questions/4938612/how-do-i-print-the-pointer-value-of-a-go-object-what-does-the-pointer-value-mea#4963935

func TestSliceAssignment(t *testing.T) {
	var blocks1 []chain.Block
	var blocks2 []chain.Block
	blocks1 = createTestBlocks(10)
	blocks2 = createTestBlocks(20)

	//block1 and block2 are different slices on different arrays, so each should have different address
	assert.NotEqual(t, &blocks1, &blocks2)
	//block1 and block2 should point to different arrays, so address of first item should be different
	assert.NotEqual(t, &blocks1[0], &blocks2[0])

	assert.Equal(t, 10, len(blocks1))
	assert.Equal(t, 20, len(blocks2))
	assert.Equal(t, 0, blocks1[0].Nonce)
	assert.Equal(t, 0, blocks2[0].Nonce)
	blocks1[0].Nonce = 10
	blocks2[0].Nonce = 20
	assert.Equal(t, 10, blocks1[0].Nonce)
	assert.Equal(t, 20, blocks2[0].Nonce)

	blocks1 = blocks2
	assert.Equal(t, 20, len(blocks1))
	assert.Equal(t, 20, len(blocks2))
	assert.Equal(t, 20, blocks1[0].Nonce)
	assert.Equal(t, 20, blocks2[0].Nonce)

	//block1 and block2 are slice objects, so apparently even if pointing to exact same data, the address of slice is different
	assert.NotEqual(t, &blocks1, &blocks2)
	//block1 and block2 should now point to same arrays, so address of first item (and all others) should be same
	assert.Equal(t, &blocks1[0], &blocks2[0])
}

func createTestBlocks(count int) []chain.Block {
	var blocks []chain.Block
	for i := 0 ; i < count ; i++ {
		block := chain.Block{0, "", "", time.Now(), "", nil, 0, 0}
		blocks = append(blocks, block)
	}
	return blocks
}

//test below seems to indicate byte(int) always produces an unsigned (positive) integer
func TestIntBytes(t *testing.T) {
	i := 1
	println("i=1:", i)
	b := byte(i)
	println("b=1:", b)

	i = -1
	println("i=-1:", i)
	b = byte(i)
	println("b=-1:", b)

	i = 128
	println("i=128:", i)
	b = byte(i)
	println("b=128:", b)

	i = 255
	println("i=255:", i)
	b = byte(i)
	println("b=255:", b)

	i = 256
	println("i=256:", i)
	b = byte(i)
	println("b=256:", b)
}

