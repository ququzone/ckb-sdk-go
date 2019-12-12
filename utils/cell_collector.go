package utils

import (
	"context"

	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/types"
)

type CellCollector struct {
	Client     rpc.Client
	LockScript *types.Script
	Capacity   uint64
	TypeScript *types.Script
	UseIndex   bool
}

func NewCellCollector(client rpc.Client, lockScript *types.Script, capacity uint64) *CellCollector {
	return &CellCollector{
		Client:     client,
		LockScript: lockScript,
		Capacity:   capacity,
		UseIndex:   false,
	}
}

func (collector *CellCollector) Collect() ([]*types.Cell, uint64, error) {
	lockHash, err := collector.LockScript.Hash()
	if err != nil {
		return nil, 0, err
	}
	if collector.UseIndex {
		return collectFromIndex(collector.Client, lockHash, collector.Capacity, collector.TypeScript)
	}
	return collectFromBlocks(collector.Client, lockHash, collector.Capacity, collector.TypeScript)
}

func collectFromIndex(client rpc.Client, lockHash types.Hash, capacity uint64, typeScript *types.Script) ([]*types.Cell, uint64, error) {
	var result []*types.Cell
	var start uint
	var total uint64
	var stop bool
	for {
		cells, err := client.GetLiveCellsByLockHash(context.Background(), lockHash, start, 50, false)
		if err != nil {
			return nil, 0, err
		}
		for i := 0; i < len(cells); i++ {
			cell := cells[i]
			if typeScript != nil {
				if !typeScript.Equals(cell.CellOutput.Type) {
					continue
				}
			} else {
				if cell.CellOutput.Type != nil {
					continue
				}
			}
			total += cell.CellOutput.Capacity
			result = append(result, &types.Cell{
				// set blockhash to empty
				BlockHash: types.Hash{},
				Capacity:  cell.CellOutput.Capacity,
				Lock:      cell.CellOutput.Lock,
				OutPoint: &types.OutPoint{
					TxHash: cell.CreatedBy.TxHash,
					Index:  cell.CreatedBy.Index,
				},
				Type: cell.CellOutput.Type,
			})
			if total >= capacity {
				stop = true
				break
			}
		}
		if stop || len(cells) < 50 {
			break
		}
	}
	return result, total, nil
}

func collectFromBlocks(client rpc.Client, lockHash types.Hash, capacity uint64, typeScript *types.Script) ([]*types.Cell, uint64, error) {
	header, err := client.GetTipHeader(context.Background())
	if err != nil {
		return nil, 0, err
	}
	var result []*types.Cell
	var start uint64
	var total uint64
	var stop bool
	for {
		end := start + 100
		if end > header.Number {
			end = header.Number
			stop = true
		}
		cells, err := client.GetCellsByLockHash(context.Background(), lockHash, start, end)
		if err != nil {
			return nil, 0, err
		}
		for i := 0; i < len(cells); i++ {
			if typeScript != nil {
				if !typeScript.Equals(cells[i].Type) {
					continue
				}
			} else {
				if cells[i].Type != nil {
					continue
				}
			}
			result = append(result, cells[i])
			total += cells[i].Capacity
			if total >= capacity {
				stop = true
				break
			}
		}
		if stop {
			break
		}
		start += 100
	}
	return result, total, nil
}
