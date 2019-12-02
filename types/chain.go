package types

type Epoch struct {
	CompactTarget uint64 `json:"compact_target"`
	Length        uint64 `json:"length"`
	Number        uint64 `json:"number"`
	StartNumber   uint64 `json:"start_number"`
}

type Header struct {
	CompactTarget    uint64 `json:"compact_target"`
	Dao              Hash   `json:"dao"`
	Epoch            uint64 `json:"epoch"`
	Hash             Hash   `json:"hash"`
	Nonce            string `json:"nonce"`
	Number           uint64 `json:"number"`
	ParentHash       Hash   `json:"parent_hash"`
	ProposalsHash    Hash   `json:"proposals_hash"`
	Timestamp        uint64 `json:"timestamp"`
	TransactionsRoot Hash   `json:"transactions_root"`
	UnclesHash       Hash   `json:"uncles_hash"`
	Version          uint   `json:"version"`
}
