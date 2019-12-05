package address

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/ququzone/ckb-sdk-go/types"
)

func TestGenerateAddress(t *testing.T) {
	script := &types.Script{
		CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
		HashType: types.HashTypeType,
		Args:     common.FromHex("0xedcda9513fa030ce4308e29245a22c022d0443bb"),
	}

	mnAddress, err := GenerateAddress(Mainnet, script)

	assert.Nil(t, err)
	assert.Equal(t, "ckb1qyqwmndf2yl6qvxwgvyw9yj95gkqytgygwasshh9m8", mnAddress)

	tnAddress, err := GenerateAddress(Testnet, script)
	assert.Nil(t, err)
	assert.Equal(t, "ckt1qyqwmndf2yl6qvxwgvyw9yj95gkqytgygwasdjf6hm", tnAddress)
}
