package blockstream_models

type TransactionResponse struct {
	Txid     string   `json:"txid"`
	Version  int      `json:"version"`
	Locktime int      `json:"locktime"`
	Vin      []Input  `json:"vin"`
	Vout     []Output `json:"vout"`
	Size     int      `json:"size"`
	Weight   int      `json:"weight"`
	Fee      int      `json:"fee"`
	Status   Status   `json:"status"`
}

type Input struct {
	Txid    string  `json:"txid"`
	Vout    int     `json:"vout"`
	Prevout Prevout `json:"prevout"`
}

type Prevout struct {
	Scriptpubkey     string `json:"scriptpubkey"`
	ScriptpubkeyType string `json:"scriptpubkey_type"`
	Value            int64  `json:"value"`
}

type Output struct {
	Scriptpubkey        string `json:"scriptpubkey"`
	ScriptpubkeyAddress string `json:"scriptpubkey_address,omitempty"`
	Value               int64  `json:"value"`
}

type Status struct {
	Confirmed   bool    `json:"confirmed"`
	BlockHeight *int    `json:"block_height,omitempty"`
	BlockHash   *string `json:"block_hash,omitempty"`
	BlockTime   *int64  `json:"block_time,omitempty"`
}

type AddressTransactionsResponse []TransactionResponse

type StandardizedTransaction struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Token   string `json:"token"`
	Amount  string `json:"amount"`
	Value   string `json:"value"`
	Address string `json:"address"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Fee     string `json:"fee"`
	Hash    string `json:"hash"`
}

type StandardizedTransactionsResponse struct {
	Transactions []StandardizedTransaction `json:"transactions"`
	TotalCount   int                       `json:"totalCount"`
}
