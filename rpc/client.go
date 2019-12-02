package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ququzone/ckb-sdk-go/types"
)

// Client for the Nervos RPC API.
type Client interface {
	// GetTipBlockNumber returns the number of blocks in the longest blockchain.
	GetTipBlockNumber(ctx context.Context) (uint64, error)

	// GetTipHeader returns the information about the tip header of the longest.
	GetTipHeader(ctx context.Context) (*types.Header, error)

	// GetCurrentEpoch returns the information about the current epoch.
	GetCurrentEpoch(ctx context.Context) (*types.Epoch, error)

	// Close close client
	Close()
}

type client struct {
	c *rpc.Client
}

func Dial(url string) (Client, error) {
	return DialContext(context.Background(), url)
}

func DialContext(ctx context.Context, url string) (Client, error) {
	c, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

func NewClient(c *rpc.Client) Client {
	return &client{c}
}

func (cli *client) Close() {
	cli.c.Close()
}

// Chain RPC

func (cli *client) GetTipBlockNumber(ctx context.Context) (uint64, error) {
	var num hexutil.Uint64
	err := cli.c.CallContext(ctx, &num, "get_tip_block_number")
	if err != nil {
		return 0, err
	}
	return uint64(num), err
}

func (cli *client) GetCurrentEpoch(ctx context.Context) (*types.Epoch, error) {
	var result epoch
	err := cli.c.CallContext(ctx, &result, "get_current_epoch")
	if err != nil {
		return nil, err
	}
	return &types.Epoch{
		CompactTarget: uint64(result.CompactTarget),
		Length:        uint64(result.Length),
		Number:        uint64(result.Number),
		StartNumber:   uint64(result.StartNumber),
	}, err
}

func (cli *client) GetTipHeader(ctx context.Context) (*types.Header, error) {
	var result header
	err := cli.c.CallContext(ctx, &result, "get_tip_header")
	if err != nil {
		return nil, err
	}
	return &types.Header{
		CompactTarget: uint64(result.CompactTarget),
		Dao: result.Dao,
		Epoch: uint64(result.Epoch),
		Hash: result.Hash,
		Nonce: result.Nonce,
		Number: uint64(result.Number),
		ParentHash: result.ParentHash,
		ProposalsHash: result.ProposalsHash,
		Timestamp: uint64(result.Timestamp),
		TransactionsRoot: result.TransactionsRoot,
		UnclesHash: result.UnclesHash,
		Version: uint(result.Version),
	}, err
}
