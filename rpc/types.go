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

type inTransaction struct {
	Version     hexutil.Uint    `json:"version"`
	CellDeps    []cellDep       `json:"cell_deps"`
	HeaderDeps  []types.Hash    `json:"header_deps"`
	Inputs      []cellInput     `json:"inputs"`
	Outputs     []cellOutput    `json:"outputs"`
	OutputsData []hexutil.Bytes `json:"outputs_data"`
	Witnesses   []hexutil.Bytes `json:"witnesses"`
}

type uncleBlock struct {
	Header    header         `json:"header"`
	Proposals []hexutil.Uint `json:"proposals"`
}

type block struct {
	Header       header         `json:"header"`
	Proposals    []hexutil.Uint `json:"proposals"`
	Transactions []transaction  `json:"transactions"`
	Uncles       []uncleBlock   `json:"uncles"`
}

type cell struct {
	BlockHash types.Hash  `json:"block_hash"`
	Capacity  hexutil.Big `json:"capacity"`
	Lock      *script     `json:"lock"`
	OutPoint  *outPoint   `json:"out_point"`
	Type      *script     `json:"type"`
}

type cellData struct {
	Content hexutil.Bytes `json:"content"`
	Hash    types.Hash    `json:"hash"`
}

type cellInfo struct {
	Data   *cellData  `json:"data"`
	Output cellOutput `json:"output"`
}

type cellWithStatus struct {
	Cell   cellInfo `json:"cell"`
	Status string   `json:"status"`
}

type transactionWithStatus struct {
	Transaction transaction `json:"transaction"`
	TxStatus    struct {
		BlockHash types.Hash              `json:"block_hash"`
		Status    types.TransactionStatus `json:"status"`
	} `json:"tx_status"`
}

type blockReward struct {
	Primary        hexutil.Big `json:"primary"`
	ProposalReward hexutil.Big `json:"proposal_reward"`
	Secondary      hexutil.Big `json:"secondary"`
	Total          hexutil.Big `json:"total"`
	TxFee          hexutil.Big `json:"tx_fee"`
}

type dryRunTransactionResult struct {
	Cycles hexutil.Uint64 `json:"cycles"`
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

func toTransaction(tx transaction) *types.Transaction {
	return &types.Transaction{
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

func toTransactions(transactions []transaction) []*types.Transaction {
	result := make([]*types.Transaction, len(transactions))
	for i := 0; i < len(transactions); i++ {
		result[i] = toTransaction(transactions[i])
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

func toUncles(uncles []uncleBlock) []*types.UncleBlock {
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

func toCells(cells []cell) []*types.Cell {
	result := make([]*types.Cell, len(cells))
	for i := 0; i < len(cells); i++ {
		cell := cells[i]
		result[i] = &types.Cell{
			BlockHash: cell.BlockHash,
			Capacity:  (*big.Int)(&cell.Capacity),
			Lock: &types.Script{
				CodeHash: cell.Lock.CodeHash,
				HashType: cell.Lock.HashType,
				Args:     cell.Lock.Args,
			},
			OutPoint: &types.OutPoint{
				TxHash: cell.OutPoint.TxHash,
				Index:  uint64(cell.OutPoint.Index),
			},
		}
		if cell.Type != nil {
			result[i].Type = &types.Script{
				CodeHash: cell.Type.CodeHash,
				HashType: cell.Type.HashType,
				Args:     cell.Type.Args,
			}
		}
	}
	return result
}

func toCellWithStatus(status cellWithStatus) *types.CellWithStatus {
	result := &types.CellWithStatus{
		Cell: &types.CellInfo{
			Output: &types.CellOutput{
				Capacity: (*big.Int)(&status.Cell.Output.Capacity),
				Lock: &types.Script{
					CodeHash: status.Cell.Output.Lock.CodeHash,
					HashType: status.Cell.Output.Lock.HashType,
					Args:     status.Cell.Output.Lock.Args,
				},
			},
		},
		Status: status.Status,
	}

	if status.Cell.Data != nil {
		result.Cell.Data = &types.CellData{
			Content: status.Cell.Data.Content,
			Hash:    status.Cell.Data.Hash,
		}
	}

	if status.Cell.Output.Type != nil {
		result.Cell.Output.Type = &types.Script{
			CodeHash: status.Cell.Output.Type.CodeHash,
			HashType: status.Cell.Output.Type.HashType,
			Args:     status.Cell.Output.Type.Args,
		}
	}

	return result
}

func fromCellDeps(deps []*types.CellDep) []cellDep {
	result := make([]cellDep, len(deps))
	for i := 0; i < len(deps); i++ {
		dep := deps[i]
		result[i] = cellDep{
			OutPoint: outPoint{
				TxHash: dep.OutPoint.TxHash,
				Index:  hexutil.Uint64(dep.OutPoint.Index),
			},
			DepType: dep.DepType,
		}
	}
	return result
}

func fromInputs(inputs []*types.CellInput) []cellInput {
	result := make([]cellInput, len(inputs))
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		result[i] = cellInput{
			Since: hexutil.Uint64(input.Since),
			PreviousOutput: outPoint{
				TxHash: input.PreviousOutput.TxHash,
				Index:  hexutil.Uint64(input.PreviousOutput.Index),
			},
		}
	}
	return result
}

func fromOutputs(outputs []*types.CellOutput) []cellOutput {
	result := make([]cellOutput, len(outputs))
	for i := 0; i < len(outputs); i++ {
		output := outputs[i]
		result[i] = cellOutput{
			Capacity: (hexutil.Big)(*output.Capacity),
			Lock: &script{
				CodeHash: output.Lock.CodeHash,
				HashType: output.Lock.HashType,
				Args:     output.Lock.Args,
			},
		}
		if output.Type != nil {
			result[i].Type = &script{
				CodeHash: output.Type.CodeHash,
				HashType: output.Type.HashType,
				Args:     output.Type.Args,
			}
		}
	}
	return result
}

func fromBytesArray(bytes [][]byte) []hexutil.Bytes {
	result := make([]hexutil.Bytes, len(bytes))
	for i, data := range bytes {
		result[i] = data
	}
	return result
}

func fromTransaction(tx *types.Transaction) inTransaction {
	result := inTransaction{
		Version:     hexutil.Uint(tx.Version),
		HeaderDeps:  tx.HeaderDeps,
		CellDeps:    fromCellDeps(tx.CellDeps),
		Inputs:      fromInputs(tx.Inputs),
		Outputs:     fromOutputs(tx.Outputs),
		OutputsData: fromBytesArray(tx.OutputsData),
		Witnesses:   fromBytesArray(tx.Witnesses),
	}
	return result
}
