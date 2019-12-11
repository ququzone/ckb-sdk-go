package utils

import (
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/types"
)

func CollectCells(client rpc.Client, script *types.Script, capacity uint64, useIndex bool) ([]*types.Cell, error) {
	// TODO
	return nil, nil
}
