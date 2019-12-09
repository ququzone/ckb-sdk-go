package transaction

import (
	"encoding/binary"
	"errors"

	"github.com/ququzone/ckb-sdk-go/crypto"
	"github.com/ququzone/ckb-sdk-go/crypto/blake2b"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

var EmptyWitnessArg = &types.WitnessArgs{
	Lock:       make([]byte, 65),
	InputType:  nil,
	OutputType: nil,
}

func NewSecp256k1SingleSigTx(scripts *utils.SystemScripts) *types.Transaction {
	return &types.Transaction{
		Version:    0,
		HeaderDeps: []types.Hash{},
		CellDeps: []*types.CellDep{
			{
				OutPoint: scripts.SecpCell.OutPoint,
				DepType:  types.DepTypeDepGroup,
			},
		},
	}
}

func AddInputsForTransaction(transaction *types.Transaction, cells []*types.Cell) ([]int, *types.WitnessArgs, error) {
	if len(cells) == 0 {
		return nil, nil, errors.New("input cells empty")
	}
	group := make([]int, len(cells))
	start := len(transaction.Witnesses)
	for i := 0; i < len(cells); i++ {
		cell := cells[i]
		input := &types.CellInput{
			Since: 0,
			PreviousOutput: &types.OutPoint{
				TxHash: cell.OutPoint.TxHash,
				Index:  cell.OutPoint.Index,
			},
		}
		transaction.Inputs = append(transaction.Inputs, input)
		transaction.Witnesses = append(transaction.Witnesses, []byte{})
		group[i] = start + i
	}
	return group, EmptyWitnessArg, nil
}

func SingleSignTransaction(transaction *types.Transaction, group []int, witnessArgs *types.WitnessArgs, key crypto.Key) error {
	data, err := witnessArgs.Serialize()
	if err != nil {
		return err
	}
	length := make([]byte, 8)
	binary.LittleEndian.PutUint64(length, uint64(len(data)))

	hash, err := transaction.ComputeHash()
	if err != nil {
		return err
	}

	message := append(hash.Bytes(), length...)
	message = append(message, data...)

	if len(group) > 1 {
		for i := 1; i < len(group); i++ {
			var data []byte
			length := make([]byte, 8)
			binary.LittleEndian.PutUint64(length, uint64(len(data)))
			message = append(message, length...)
			message = append(message, data...)
		}
	}

	message, err = blake2b.Blake256(message)
	if err != nil {
		return err
	}

	signd, err := key.Sign(message)
	if err != nil {
		return err
	}

	wa := &types.WitnessArgs{
		Lock: signd,
	}
	wab, err := wa.Serialize()
	if err != nil {
		return err
	}

	transaction.Witnesses[group[0]] = wab

	return nil
}
