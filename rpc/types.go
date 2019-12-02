package rpc

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ququzone/ckb-sdk-go/types"
)

type epoch struct {
	CompactTarget hexutil.Uint64 `json:"compact_target"`
	Length        hexutil.Uint64 `json:"length"`
	Number        hexutil.Uint64 `json:"number"`
	StartNumber   hexutil.Uint64 `json:"start_number"`
}

type header struct {
	CompactTarget    hexutil.Uint64 `json:"compact_target"`
	Dao              types.Hash     `json:"dao"`
	Epoch            hexutil.Uint64 `json:"epoch"`
	Hash             types.Hash     `json:"hash"`
	Nonce            string         `json:"nonce"`
	Number           hexutil.Uint64 `json:"number"`
	ParentHash       types.Hash     `json:"parent_hash"`
	ProposalsHash    types.Hash     `json:"proposals_hash"`
	Timestamp        hexutil.Uint64 `json:"timestamp"`
	TransactionsRoot types.Hash     `json:"transactions_root"`
	UnclesHash       types.Hash     `json:"uncles_hash"`
	Version          hexutil.Uint   `json:"version"`
}
