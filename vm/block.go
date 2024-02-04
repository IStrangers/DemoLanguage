package vm

type BlockType int

const (
	_ BlockType = iota
	BlockScope
	BlockLoop
	BlockSwitch
	BlockIterator
	BlockTry
)

type Block struct {
	outer        *Block
	blockType    BlockType
	breaks       []int
	continueBase int
	continues    []int
}
