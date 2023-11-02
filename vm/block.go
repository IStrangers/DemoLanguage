package vm

type BlockType int

const (
	_ BlockType = iota
)

type Block struct {
	outer     *Block
	blockType BlockType
	label     string
}
