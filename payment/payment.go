package payment

import (
	"context"
	"fmt"

	"github.com/ququzone/ckb-sdk-go/address"
	"github.com/ququzone/ckb-sdk-go/crypto"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

type Payment struct {
	From        *types.Script
	To          *types.Script
	Amount      uint64
	Fee         uint64
	group       []int
	witnessArgs *types.WitnessArgs
	tx          *types.Transaction
}

func NewPayment(from, to string, amount, fee uint64) (*Payment, error) {
	fromAddress, err := address.Parse(from)
	if err != nil {
		return nil, fmt.Errorf("parse from address %s error: %v", from, err)
	}
	toAddress, err := address.Parse(to)
	if err != nil {
		return nil, fmt.Errorf("parse to address %s error: %v", to, err)
	}

	if fromAddress.Mode != toAddress.Mode {
		return nil, fmt.Errorf("from address and to address with diffrent network: %v:%v", fromAddress.Mode, toAddress.Mode)
	}

	return &Payment{
		From:   fromAddress.Script,
		To:     toAddress.Script,
		Amount: amount,
		Fee:    fee,
	}, nil
}

func (p *Payment) GenerateTx(client rpc.Client) (*types.Transaction, error) {
	collector := utils.NewCellCollector(client, p.From, p.Amount+p.Fee)

	cells, total, err := collector.Collect()
	if err != nil {
		return nil, fmt.Errorf("collect cell error: %v", err)
	}

	if total < p.Amount+p.Fee {
		return nil, fmt.Errorf("insufficient balance: %d", total)
	}

	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		return nil, fmt.Errorf("load system script error: %v", err)
	}

	tx := transaction.NewSecp256k1SingleSigTx(systemScripts)
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: p.Amount,
		Lock:     p.To,
	})
	tx.OutputsData = [][]byte{{}}

	if total-p.Amount-p.Fee > 0 {
		tx.Outputs = append(tx.Outputs, &types.CellOutput{
			Capacity: total - p.Amount - p.Fee,
			Lock:     p.From,
		})
		tx.OutputsData = [][]byte{{}, {}}
	}

	group, witnessArgs, err := transaction.AddInputsForTransaction(tx, cells)
	if err != nil {
		return nil, fmt.Errorf("add inputs to transaction error: %v", err)
	}

	p.group = group
	p.witnessArgs = witnessArgs
	p.tx = tx
	return tx, err
}

func (p *Payment) Sign(key crypto.Key) (*types.Transaction, error) {
	err := transaction.SingleSignTransaction(p.tx, p.group, p.witnessArgs, key)
	if err != nil {
		return nil, fmt.Errorf("sign transaction error: %v", err)
	}

	return p.tx, err
}

func (p *Payment) Send(client rpc.Client) (*types.Hash, error) {
	return client.SendTransaction(context.Background(), p.tx)
}
