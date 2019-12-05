package address

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ququzone/ckb-sdk-go/crypto/bech32"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
)

type Mode string

const (
	Mainnet Mode = "ckb"
	Testnet Mode = "ckt"

	SHORT_FORMAT                 = "01"
	FULL_DATA_FORMAT             = "02"
	FULL_TYPE_FORMAT             = "04"
	CODE_HASH_INDEX_SINGLESIG    = "00"
	CODE_HASH_INDEX_MULTISIG_SIG = "01"
)

func GenerateAddress(mode Mode, script *types.Script) (string, error) {
	if script.HashType == types.HashTypeType && len(script.Args) == 20 {
		if transaction.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH == script.CodeHash.String() {
			// generate_short_payload_singlesig_address
			payload := SHORT_FORMAT + CODE_HASH_INDEX_SINGLESIG + hex.EncodeToString(script.Args)
			data, err := bech32.ConvertBits(common.FromHex(payload), 8, 5, true)
			if err != nil {
				return "", err
			}
			return bech32.Encode((string)(mode), data)
		} else if transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH == script.CodeHash.String() {
			// generate_short_payload_multisig_address
			payload := SHORT_FORMAT + CODE_HASH_INDEX_MULTISIG_SIG + hex.EncodeToString(script.Args)
			data, err := bech32.ConvertBits(common.FromHex(payload), 8, 5, true)
			if err != nil {
				return "", err
			}
			return bech32.Encode((string)(mode), data)
		} else {
			// generate_full_payload_address
			return generateFullPayloadAddress(FULL_TYPE_FORMAT, mode, script)
		}
	}

	return generateFullPayloadAddress(FULL_DATA_FORMAT, mode, script)
}

func generateFullPayloadAddress(hashType string, mode Mode, script *types.Script) (string, error) {
	payload := hashType + hex.EncodeToString(script.CodeHash.Bytes()) + hex.EncodeToString(script.Args)
	data, err := bech32.ConvertBits(common.FromHex(payload), 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode((string)(mode), data)
}
