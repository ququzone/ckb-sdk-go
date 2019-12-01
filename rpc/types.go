package rpc

import "github.com/ethereum/go-ethereum/common/hexutil"

type epoch struct {
	CompactTarget hexutil.Uint64 `json:"compact_target"`
	Length        hexutil.Uint64 `json:"length"`
	Number        hexutil.Uint64 `json:"number"`
	StartNumber   hexutil.Uint64 `json:"start_number"`
}
