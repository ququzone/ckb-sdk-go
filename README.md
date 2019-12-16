# CKB SDK Golang

[![License](https://img.shields.io/badge/license-MIT-green)](https://github.com/ququzone/ckb-sdk-go/blob/master/LICENSE)
[![Go version](https://img.shields.io/badge/go-1.11.5-blue.svg)](https://github.com/moovweb/gvm)
[![Telegram Group](https://cdn.rawgit.com/Patrolavia/telegram-badge/8fe3382b/chat.svg)](https://t.me/nervos_ckb_dev)

Golang SDK for Nervos [CKB](https://github.com/nervosnetwork/ckb).

## Get started

### Minimum requirements

| Components | Version | Description |
|----------|-------------|-------------|
| [Golang](https://golang.org) | &ge; 1.11.5 | Go programming language |

### Install

```bash
go get -v github.com/ququzone/ckb-sdk-go
```

## Basic Usages

### 1. Single input send transaction

```go
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}

	key, err := secp256k1.HexToKey(PRIVATE_KEY)
	if err != nil {
		log.Fatalf("import private key error: %v", err)
	}

	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}

	toAddress, _ := hex.DecodeString("bf3e92da4911fa5f620e7b1fd27c2d0ddd0de744")
	changeScript, _ := key.Script(systemScripts)

	tx := transaction.NewSecp256k1SingleSigTx(systemScripts)
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: 200000000000,
		Lock: &types.Script{
			CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			HashType: types.HashTypeType,
			Args:     toAddress,
		},
	})
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: 199999998000,
		Lock: &types.Script{
			CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			HashType: types.HashTypeType,
			Args:     changeScript.Args,
		},
	})
	tx.OutputsData = [][]byte{{}, {}}

	group, witnessArgs, err := transaction.AddInputsForTransaction(tx, []*types.Cell{
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x8e6d818c6e07e6cbd9fca51294030494ee23dc388d7f5276ba50b938d02cc015"),
				Index:  1,
			},
		},
	})
	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	err = transaction.SingleSignTransaction(tx, group, witnessArgs, key)
	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send transaction error: %v", err)
	}

	fmt.Println(hash.String())
}
```

### 2. Multiple inputs send transaction

```go
package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}

	keyA, err := secp256k1.HexToKey(PRIVATE_KEY_A)
	if err != nil {
		log.Fatalf("import private key error: %v", err)
	}

	keyB, err := secp256k1.HexToKey(PRIVATE_KEY_B)
	if err != nil {
		log.Fatalf("import private key error: %v", err)
	}

	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}

	toAddress, _ := hex.DecodeString("f96b6700df60fd6d84a2e17a5c5e4f603a5eec5d")

	tx := transaction.NewSecp256k1SingleSigTx(systemScripts)
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: 499999996000,
		Lock: &types.Script{
			CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
			HashType: types.HashTypeType,
			Args:     toAddress,
		},
	})
	tx.OutputsData = [][]byte{{}}

	groupB, witnessArgsB, err := transaction.AddInputsForTransaction(tx, []*types.Cell{
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xf56d73acbe235889e726366aa4fa09b3f0b51138c294645bb30912fb872837a5"),
				Index:  0,
			},
		},
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x8e6d818c6e07e6cbd9fca51294030494ee23dc388d7f5276ba50b938d02cc015"),
				Index:  0,
			},
		},
	})
	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	groupA, witnessArgsA, err := transaction.AddInputsForTransaction(tx, []*types.Cell{
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xf56d73acbe235889e726366aa4fa09b3f0b51138c294645bb30912fb872837a5"),
				Index:  1,
			},
		},
	})
	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	err = transaction.SingleSignTransaction(tx, groupB, witnessArgsB, keyB)
	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	err = transaction.SingleSignTransaction(tx, groupA, witnessArgsA, keyA)
	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send transaction error: %v", err)
	}

	fmt.Println(hash.String())
}
```

### 3. Multiple inputs and multisig transaction

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ququzone/ckb-sdk-go/address"
	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}

	keyA, _ := secp256k1.HexToKey(PRIVATE_KEY_A)
	keyB, _ := secp256k1.HexToKey(PRIVATE_KEY_B)
	keyC, _ := secp256k1.HexToKey(PRIVATE_KEY_C)

	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}

	to, _ := address.Parse("ckt1qyqwmndf2yl6qvxwgvyw9yj95gkqytgygwasdjf6hm")
	changeScript, changeSerialize, err := address.GenerateSecp256k1MultisigScript(0, 2, [][]byte{
		keyA.PubKey(),
		keyB.PubKey(),
		keyC.PubKey(),
	})

	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}

	tx := transaction.NewSecp256k1MultiSigTx(systemScripts)
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: 300000000000,
		Lock:     to.Script,
	})
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: 399999994000,
		Lock:     changeScript,
	})
	tx.OutputsData = [][]byte{{}, {}}

	group, witnessArgs, err := transaction.AddInputsForTransaction(tx, []*types.Cell{
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xccb33a76b5322ff2841511ef10606b6bb207f6eef5a687f14f8c7fa5da8a7cb2"),
				Index:  0,
			},
		},
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0x06a49393423c1be0a48d422fa60951bdb847d56753915f321c26906a6ba1dd8a"),
				Index:  0,
			},
		},
	})
	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	err = transaction.MultiSignTransaction(tx, group, witnessArgs, changeSerialize, keyA, keyB)
	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send transaction error: %v", err)
	}

	fmt.Println(hash.String())
}
```

### 3. Multiple inputs and hybirdsig transaction

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ququzone/ckb-sdk-go/address"
	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}

	key, _ := secp256k1.HexToKey(PRIVATE_KEY)

	keyA, _ := secp256k1.HexToKey(PRIVATE_KEY_A)
	keyB, _ := secp256k1.HexToKey(PRIVATE_KEY_B)
	keyC, _ := secp256k1.HexToKey(PRIVATE_KEY_C)

	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}

	to, _ := address.Parse("ckt1qyqt705jmfy3r7jlvg88k87j0sksmhgduazq7x5l8k")
	changeScript, changeSerialize, err := address.GenerateSecp256k1MultisigScript(0, 2, [][]byte{
		keyA.PubKey(),
		keyB.PubKey(),
		keyC.PubKey(),
	})

	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}

	tx := transaction.NewSecp256k1HybirdSigTx(systemScripts)
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: 500000000000,
		Lock:     to.Script,
	})
	tx.Outputs = append(tx.Outputs, &types.CellOutput{
		Capacity: 199999992000,
		Lock:     changeScript,
	})
	tx.OutputsData = [][]byte{{}, {}}

	group, witnessArgs, err := transaction.AddInputsForTransaction(tx, []*types.Cell{
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xcb905a3b304b23200225def794c4ce165d93eead77197724680b4ec067b43803"),
				Index:  0,
			},
		},
	})
	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	group1, witnessArgs1, err := transaction.AddInputsForTransaction(tx, []*types.Cell{
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xcb905a3b304b23200225def794c4ce165d93eead77197724680b4ec067b43803"),
				Index:  1,
			},
		},
	})
	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	err = transaction.SingleSignTransaction(tx, group, witnessArgs, key)
	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	err = transaction.MultiSignTransaction(tx, group1, witnessArgs1, changeSerialize, keyA, keyB)
	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	hash, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send transaction error: %v", err)
	}

	fmt.Println(hash.String())
}
```

### 5. Collect cells

```go
package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}

	args, _ := hex.DecodeString("6407c2ef9bd96e8e14ac4cd15d860e9331802172")

	collector := utils.NewCellCollector(client, &types.Script{
		CodeHash: types.HexToHash("0x9bd7e06f3ecf4be0f2fcd2188b23f1b9fcc88e5d4b65a8637b17723bbda3cce8"),
		HashType: types.HashTypeType,
		Args:     args,
	}, 10000000000000000)

	// default collect null type script
	fmt.Println(collector.Collect())

	// collect by type script
	collector.TypeScript = &types.Script{
		CodeHash: types.HexToHash("0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e"),
		HashType: types.HashTypeType,
		Args:     []byte{},
	}

	cells, total, err := collector.Collect()

	fmt.Println(total)
	fmt.Println(err)
	fmt.Println(cells)
}
```

### 6. Payment

```go
package main

import (
	"fmt"
	"log"

	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/payment"
	"github.com/ququzone/ckb-sdk-go/rpc"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}

	key, err := secp256k1.HexToKey(PRIVATE_KEY)
	if err != nil {
		log.Fatalf("import private key error: %v", err)
	}

	pay, err := payment.NewPayment("ckt1qyqwmndf2yl6qvxwgvyw9yj95gkqytgygwasdjf6hm",
		"ckt1qyqt705jmfy3r7jlvg88k87j0sksmhgduazq7x5l8k", 100000000000, 1000)
	if err != nil {
		log.Fatalf("create payment error: %v", err)
	}

	_, err = pay.GenerateTx(client)
	if err != nil {
		log.Fatalf("create transaction error: %v", err)
	}

	_, err = pay.Sign(key)
	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	hash, err := pay.Send(client)

	fmt.Println(hash)
}
```

### 7. Dao deposit

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ququzone/ckb-sdk-go/address"
	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/dao"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}

	key, err := secp256k1.HexToKey(PRIVATE_KEY)
	if err != nil {
		log.Fatalf("import private key error: %v", err)
	}

	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}

	deposit := dao.NewDeposit(systemScripts, false)

	to, _ := address.Parse("ckt1qyqwmndf2yl6qvxwgvyw9yj95gkqytgygwasdjf6hm")
	change, _ := key.Script(systemScripts)

	err = deposit.AddDaoOutput(systemScripts, to.Script, 400000000000)
	if err != nil {
		log.Fatalf("add dao output error: %v", err)
	}
	err = deposit.AddOutput(change, 99999997000)
	if err != nil {
		log.Fatalf("add output error: %v", err)
	}

	group, witnessArgs, err := transaction.AddInputsForTransaction(deposit.Transaction, []*types.Cell{
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xaa10f51bc6ee60e851d17e3fffefc950d6dc1d2cd77e15699c3da5e837219764"),
				Index:  1,
			},
		},
	})
	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	err = transaction.SingleSignTransaction(deposit.Transaction, group, witnessArgs, key)
	if err != nil {
		log.Fatalf("sign transaction error: %v", err)
	}

	hash, err := client.SendTransaction(context.Background(), deposit.Transaction)
	if err != nil {
		log.Fatalf("send transaction error: %v", err)
	}
	fmt.Println(hash.String())
}
```

### 8. Dao withdraw Phase1

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ququzone/ckb-sdk-go/address"
	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/dao"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}

	key, err := secp256k1.HexToKey(PRIVATE_KEY)
	if err != nil {
		log.Fatalf("import private key error: %v", err)
	}

	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}

	withdraw := dao.NewWithdrawPhase1(systemScripts, false)

	ownder, _ := address.Parse("ckt1qyqwmndf2yl6qvxwgvyw9yj95gkqytgygwasdjf6hm")
	change, _ := key.Script(systemScripts)

	index, err := withdraw.AddDaoDepositTick(client, &types.Cell{
		BlockHash: types.HexToHash("0x386bafd53bade6bf769c9b10f545e31ea744cb6ebc5f1c8178f307e8dce157a6"),
		Capacity:  400000000000,
		Lock:      ownder.Script,
		Type: &types.Script{
			CodeHash: types.HexToHash("0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e"),
			HashType: types.HashTypeType,
			Args:     []byte{},
		},
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash("0xc8cfe3d09b0a50fd2df3bd79dbadca23b7eb1f58087942d7266abea93459fce1"),
			Index:  0,
		},
	})
	if err != nil {
		log.Fatalf("add dao deposit tick error: %v", err)
	}

	err = withdraw.AddOutput(change, 99999995000)
	if err != nil {
		log.Fatalf("add output error: %v", err)
	}

	group, witnessArgs, err := transaction.AddInputsForTransaction(withdraw.Transaction, []*types.Cell{
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xc8cfe3d09b0a50fd2df3bd79dbadca23b7eb1f58087942d7266abea93459fce1"),
				Index:  1,
			},
		},
	})
	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	var groups []int
	groups = append(groups, index)
	groups = append(groups, group...)

	// sign dao input
	err = transaction.SingleSignTransaction(withdraw.Transaction, groups, witnessArgs, key)
	if err != nil {
		log.Fatalf("sign dao transaction error: %v", err)
	}

	hash, err := client.SendTransaction(context.Background(), withdraw.Transaction)
	if err != nil {
		log.Fatalf("send transaction error: %v", err)
	}
	fmt.Println(hash.String())
}
```

### 8. Dao withdraw phase2

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ququzone/ckb-sdk-go/address"
	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/dao"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
)

func main() {
	client, err := rpc.Dial("http://127.0.0.1:8114")
	if err != nil {
		log.Fatalf("create rpc client error: %v", err)
	}

	key, err := secp256k1.HexToKey(PRIVATE_KEY)
	if err != nil {
		log.Fatalf("import private key error: %v", err)
	}

	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		log.Fatalf("load system script error: %v", err)
	}

	withdraw := dao.NewWithdrawPhase2(systemScripts, false)

	ownder, _ := address.Parse("ckt1qyqwmndf2yl6qvxwgvyw9yj95gkqytgygwasdjf6hm")
	change, _ := key.Script(systemScripts)

	index, headerIndex, err := withdraw.AddDaoWithdrawTick(client, &types.Cell{
		BlockHash: types.HexToHash("0x386bafd53bade6bf769c9b10f545e31ea744cb6ebc5f1c8178f307e8dce157a6"),
		Capacity:  400000000000,
		Lock:      ownder.Script,
		Type: &types.Script{
			CodeHash: types.HexToHash("0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e"),
			HashType: types.HashTypeType,
			Args:     []byte{},
		},
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash("0xc8cfe3d09b0a50fd2df3bd79dbadca23b7eb1f58087942d7266abea93459fce1"),
			Index:  0,
		},
	}, &types.Cell{
		BlockHash: types.HexToHash("0xf0b88e5ca0397c761dc76fa2dd953f203b90c47a7c8199b45ac0d90efb044233"),
		Capacity:  400000000000,
		Lock:      ownder.Script,
		Type: &types.Script{
			CodeHash: types.HexToHash("0x82d76d1b75fe2fd9a27dfbaa65a039221a380d76c926f378d3f81cf3e7e13f2e"),
			HashType: types.HashTypeType,
			Args:     []byte{},
		},
		OutPoint: &types.OutPoint{
			TxHash: types.HexToHash("0xc72d7bffcc3302f8267fecb103f655e63e7b94b6f6e863cd6a0130ffec296684"),
			Index:  0,
		},
	})
	if err != nil {
		log.Fatalf("add dao deposit tick error: %v", err)
	}

	err = withdraw.AddOutput(change, 99999993000)
	if err != nil {
		log.Fatalf("add output error: %v", err)
	}

	group, witnessArgs, err := transaction.AddInputsForTransaction(withdraw.Transaction, []*types.Cell{
		{
			OutPoint: &types.OutPoint{
				TxHash: types.HexToHash("0xc72d7bffcc3302f8267fecb103f655e63e7b94b6f6e863cd6a0130ffec296684"),
				Index:  1,
			},
		},
	})
	if err != nil {
		log.Fatalf("add inputs to transaction error: %v", err)
	}

	var groups []int
	groups = append(groups, index)
	groups = append(groups, group...)

	// sign dao input
	err = transaction.SingleSignTransaction(withdraw.Transaction, groups, witnessArgs, key)
	if err != nil {
		log.Fatalf("sign dao transaction error: %v", err)
	}

	withdraw.Transaction.Witnesses[index] = append(withdraw.Transaction.Witnesses[index], types.SerializeUint64(headerIndex)...)

	hash, err := client.SendTransaction(context.Background(), withdraw.Transaction)
	if err != nil {
		log.Fatalf("send transaction error: %v", err)
	}
	fmt.Println(hash.String())
}
```