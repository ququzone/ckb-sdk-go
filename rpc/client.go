package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// Client for the Nervos RPC API.
type Client struct {
	c *rpc.Client
}

func Dial(url string) (*Client, error) {
	return DialContext(context.Background(), url)
}

func DialContext(ctx context.Context, url string) (*Client, error) {
	c, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

func NewClient(c *rpc.Client) *Client {
	return &Client{c}
}

func (client *Client) Close() {
	client.c.Close()
}



// Chain RPC

// GetTipBlockBumber returns the number of blocks in the longest blockchain.
func (client *Client) GetTipBlockBumber(ctx context.Context) (uint64, error) {
	var num hexutil.Uint64
	err := client.c.CallContext(ctx, &num, "get_tip_block_number")
	if err != nil {
		return 0, err
	}
	return uint64(num), err
}
