package accounts

type Account struct {
	PubKey      []byte  `json:"pubkey,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Denom       string  `json:"denom,omitempty"`
	Nonce       int64   `json:"nonce,omitempty"`
	Txn_pending bool    `json:"txn,omitempty"`
}
