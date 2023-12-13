package vm

type BlockType int

const (
	_ BlockType = iota
	BlockScope
	BlockLoop
	BlockSwitch
)

type Block struct {
	outer        *Block
	blockType    BlockType
	breaks       []int
	continueBase int
	continues    []int
}
