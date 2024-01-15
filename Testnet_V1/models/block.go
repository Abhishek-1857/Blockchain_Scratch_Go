package models

type Block struct {
	BlockHash   string
	BlockHeader Header
	MetaData    []string
}

type Header struct {
	Merkleroot      string
	Datahash        string
	Prevhash        string
	Proposeraddress string
	Timestamp       string
	Height          uint64
	TransactionCnt  int
}
