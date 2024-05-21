package blockchain

import (
	"encoding/hex"
	"fmt"
)

func (tx *Transaction) validatePrevTxns(prevTxns map[string]Transaction) error {
	for _, in := range tx.Inputs {
		if prevTxns[hex.EncodeToString(in.ID)].ID == nil {
			return fmt.Errorf("ERROR: transaction does not exist " + hex.EncodeToString(in.ID))
		}
	}

	return nil
}
