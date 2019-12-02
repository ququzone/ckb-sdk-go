package rpc

import (
	"math/big"

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
	Nonce            hexutil.Big    `json:"nonce"`
	Number           hexutil.Uint64 `json:"number"`
	ParentHash       types.Hash     `json:"parent_hash"`
	ProposalsHash    types.Hash     `json:"proposals_hash"`
	Timestamp        hexutil.Uint64 `json:"timestamp"`
	TransactionsRoot types.Hash     `json:"transactions_root"`
	UnclesHash       types.Hash     `json:"uncles_hash"`
	Version          hexutil.Uint   `json:"version"`
}

type outPoint struct {
	TxHash types.Hash     `json:"tx_hash"`
	Index  hexutil.Uint64 `json:"index"`
}

type cellDep struct {
	OutPoint outPoint      `json:"out_point"`
	DepType  types.DepType `json:"dep_type"`
}

type cellInput struct {
	Since          hexutil.Uint64 `json:"since"`
	PreviousOutput outPoint       `json:"previous_output"`
}

type script struct {
	CodeHash types.Hash           `json:"code_hash"`
	HashType types.ScriptHashType `json:"hash_type"`
	Args     hexutil.Bytes        `json:"args"`
}

type cellOutput struct {
	Capacity hexutil.Big `json:"capacity"`
	Lock     *script     `json:"lock"`
	Type     *script     `json:"type"`
}

type transaction struct {
	Version     hexutil.Uint    `json:"version"`
	Hash        types.Hash      `json:"hash"`
	CellDeps    []cellDep       `json:"cell_deps"`
	HeaderDeps  []types.Hash    `json:"header_deps"`
	Inputs      []cellInput     `json:"inputs"`
	Outputs     []cellOutput    `json:"outputs"`
	OutputsData []hexutil.Bytes `json:"outputs_data"`
	Witnesses   []hexutil.Bytes `json:"witnesses"`
}

type UncleBlock struct {
	Header    header         `json:"header"`
	Proposals []hexutil.Uint `json:"proposals"`
}

type block struct {
	Header       header         `json:"header"`
	Proposals    []hexutil.Uint `json:"proposals"`
	Transactions []transaction  `json:"transactions"`
	Uncles       []UncleBlock   `json:"uncles"`
}

func toHeader(head header) *types.Header {
	return &types.Header{
		CompactTarget:    uint64(head.CompactTarget),
		Dao:              head.Dao,
		Epoch:            uint64(head.Epoch),
		Hash:             head.Hash,
		Nonce:            (*big.Int)(&head.Nonce),
		Number:           uint64(head.Number),
		ParentHash:       head.ParentHash,
		ProposalsHash:    head.ProposalsHash,
		Timestamp:        uint64(head.Timestamp),
		TransactionsRoot: head.TransactionsRoot,
		UnclesHash:       head.UnclesHash,
		Version:          uint(head.Version),
	}
}

func toTransactions(transactions []transaction) []*types.Transaction {
	result := make([]*types.Transaction, len(transactions))
	for i := 0; i < len(transactions); i++ {
		tx := transactions[i]
		result[i] = &types.Transaction{
			Version:     uint(tx.Version),
			Hash:        tx.Hash,
			CellDeps:    toCellDeps(tx.CellDeps),
			HeaderDeps:  tx.HeaderDeps,
			Inputs:      toInputs(tx.Inputs),
			Outputs:     toOutputs(tx.Outputs),
			OutputsData: toBytesArray(tx.OutputsData),
			Witnesses:   toBytesArray(tx.Witnesses),
		}
	}
	return result
}

func toBytesArray(bytes []hexutil.Bytes) [][]byte {
	result := make([][]byte, len(bytes))
	for i, data := range bytes {
		result[i] = data
	}
	return result
}

func toOutputs(outputs []cellOutput) []*types.CellOutput {
	result := make([]*types.CellOutput, len(outputs))
	for i := 0; i < len(outputs); i++ {
		output := outputs[i]
		result[i] = &types.CellOutput{
			Capacity: (*big.Int)(&output.Capacity),
			Lock: &types.Script{
				CodeHash: output.Lock.CodeHash,
				HashType: output.Lock.HashType,
				Args:     output.Lock.Args,
			},
		}
		if output.Type != nil {
			result[i].Type = &types.Script{
				CodeHash: output.Type.CodeHash,
				HashType: output.Type.HashType,
				Args:     output.Type.Args,
			}
		}
	}
	return result
}

func toInputs(inputs []cellInput) []*types.CellInput {
	result := make([]*types.CellInput, len(inputs))
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		result[i] = &types.CellInput{
			Since: uint64(input.Since),
			PreviousOutput: &types.OutPoint{
				TxHash: input.PreviousOutput.TxHash,
				Index:  uint64(input.PreviousOutput.Index),
			},
		}
	}
	return result
}

func toCellDeps(deps []cellDep) []*types.CellDep {
	result := make([]*types.CellDep, len(deps))
	for i := 0; i < len(deps); i++ {
		dep := deps[i]
		result[i] = &types.CellDep{
			OutPoint: &types.OutPoint{
				TxHash: dep.OutPoint.TxHash,
				Index:  uint64(dep.OutPoint.Index),
			},
			DepType: dep.DepType,
		}
	}
	return result
}

func toUints(uints []hexutil.Uint) []uint {
	result := make([]uint, len(uints))
	for i, value := range uints {
		result[i] = uint(value)
	}
	return result
}

func toUncles(uncles []UncleBlock) []*types.UncleBlock {
	result := make([]*types.UncleBlock, len(uncles))
	for i := 0; i < len(uncles); i++ {
		block := uncles[i]
		result[i] = &types.UncleBlock{
			Header:    toHeader(block.Header),
			Proposals: toUints(block.Proposals),
		}
	}
	return result
}
