package dao

import (
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

type Deposit struct {
	Transaction *types.Transaction
}

func NewDeposit(scripts *utils.SystemScripts, isMultisig bool) *Deposit {
	var baseDep *types.CellDep
	if isMultisig {
		baseDep = &types.CellDep{
			OutPoint: scripts.MultiSigCell.OutPoint,
			DepType:  types.DepTypeDepGroup,
		}
	} else {
		baseDep = &types.CellDep{
			OutPoint: scripts.SecpCell.OutPoint,
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

	return &Deposit{
		Transaction: tx,
	}
}
