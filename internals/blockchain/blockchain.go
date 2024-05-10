package blockchain

type Blockchain struct {
	Blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) {
	prev := bc.Blocks[len(bc.Blocks)-1]
	bc.Blocks = append(bc.Blocks, NewBlock(prev.Hash, []byte(data)))
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{NewGenesisBlock()},
	}
}
