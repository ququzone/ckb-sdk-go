package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

type WithdrawPhase1 struct {
	Transaction *types.Transaction
}

func NewWithdrawPhase1(scripts *utils.SystemScripts, isMultisig bool) *WithdrawPhase1 {
	var baseDep *types.CellDep
	if isMultisig {
		baseDep = &types.CellDep{
			OutPoint: scripts.SecpMultiSigCell.OutPoint,
			DepType:  types.DepTypeDepGroup,
		}
	} else {
		baseDep = &types.CellDep{
			OutPoint: scripts.SecpSingleSigCell.OutPoint,
			DepType:  types.DepTypeDepGroup,
		}
	}

	tx := &types.Transaction{
		Version:    0,
		HeaderDeps: []types.Hash{},
		CellDeps: []*types.CellDep{
			baseDep,
			{
				OutPoint: scripts.DaoCell.OutPoint,
				DepType:  types.DepTypeCode,
			},
		},
	}

	return &WithdrawPhase1{
		Transaction: tx,
	}
}

func (w *WithdrawPhase1) AddDaoDepositTick(client rpc.Client, cell *types.Cell) (int, error) {
	header, err := client.GetHeader(context.Background(), cell.BlockHash)
	if err != nil {
		return 0, fmt.Errorf("get block header from address %s error: %v", cell.BlockHash.String(), err)
	}

	w.Transaction.HeaderDeps = append(w.Transaction.HeaderDeps, cell.BlockHash)
	w.Transaction.Inputs = append(w.Transaction.Inputs, &types.CellInput{
		Since: 0,
		PreviousOutput: &types.OutPoint{
			TxHash: cell.OutPoint.TxHash,
			Index:  cell.OutPoint.Index,
		},
	})
	w.Transaction.Witnesses = append(w.Transaction.Witnesses, []byte{})
	w.Transaction.Outputs = append(w.Transaction.Outputs, &types.CellOutput{
		Capacity: cell.Capacity,
		Lock:     cell.Lock,
		Type:     cell.Type,
	})
	w.Transaction.OutputsData = append(w.Transaction.OutputsData, types.SerializeUint64(header.Number))
	return len(w.Transaction.Inputs) - 1, nil
}

func (w *WithdrawPhase1) AddOutput(lock *types.Script, amount uint64) error {
	if w.Transaction == nil {
		return errors.New("must init transaction first")
	}
	w.Transaction.Outputs = append(w.Transaction.Outputs, &types.CellOutput{
		Capacity: amount,
		Lock:     lock,
	})
	w.Transaction.OutputsData = append(w.Transaction.OutputsData, []byte{})

	return nil
}
