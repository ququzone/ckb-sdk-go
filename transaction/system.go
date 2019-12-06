package transaction

import (
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/types"
)

type SystemScriptCell struct {
	CellHash types.Hash
	OutPoint *types.OutPoint
}

type SystemScripts struct {
	SecpCell     *SystemScriptCell
	MultiSigCell *SystemScriptCell
	DaoCell      *SystemScriptCell
}

func NewSystemScript(client *rpc.Client) (*SystemScripts, error) {
	// TODO
	return nil, nil
}
