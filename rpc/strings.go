package rpc

import (
	"encoding/json"

	"github.com/ququzone/ckb-sdk-go/types"
)

func TransactionString(tx *types.Transaction) (string, error) {
	itx := fromTransaction(tx)
	bytes, err := json.Marshal(itx)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
