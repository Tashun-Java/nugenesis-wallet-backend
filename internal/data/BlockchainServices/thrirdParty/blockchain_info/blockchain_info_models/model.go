package blockchain_info_models

//import "time"

type AddressInfo struct {
	Hash160       string `json:"hash160"`
	Address       string `json:"address"`
	NTx           int    `json:"n_tx"`
	NUnredeemed   int    `json:"n_unredeemed"`
	TotalReceived int64  `json:"total_received"`
	TotalSent     int64  `json:"total_sent"`
	FinalBalance  int64  `json:"final_balance"`
	Txs           []Tx   `json:"txs,omitempty"`
}

type Tx struct {
	Hash        string   `json:"hash"`
	Ver         int      `json:"ver"`
	VinSz       int      `json:"vin_sz"`
	VoutSz      int      `json:"vout_sz"`
	Size        int      `json:"size"`
	Weight      int      `json:"weight"`
	Fee         int64    `json:"fee"`
	RelayedBy   string   `json:"relayed_by"`
	LockTime    int64    `json:"lock_time"`
	TxIndex     int64    `json:"tx_index"`
	DoubleSpend bool     `json:"double_spend"`
	Time        int64    `json:"time"`
	BlockIndex  int64    `json:"block_index"`
	BlockHeight int      `json:"block_height"`
	Inputs      []Input  `json:"inputs"`
	Out         []Output `json:"out"`
}

type Input struct {
	Sequence int64    `json:"sequence"`
	Witness  string   `json:"witness"`
	Script   string   `json:"script"`
	Index    int      `json:"index"`
	PrevOut  *PrevOut `json:"prev_out"`
}

type PrevOut struct {
	Spent   bool   `json:"spent"`
	TxIndex int64  `json:"tx_index"`
	Type    int    `json:"type"`
	Addr    string `json:"addr"`
	Value   int64  `json:"value"`
	N       int    `json:"n"`
	Script  string `json:"script"`
}

type Output struct {
	Type             int   `json:"type"`
	Spent            bool  `json:"spent"`
	Value            int64 `json:"value"`
	SpendingOutpoint *struct {
		TxIndex int64 `json:"tx_index"`
		N       int   `json:"n"`
	} `json:"spending_outpoint"`
	N       int    `json:"n"`
	TxIndex int64  `json:"tx_index"`
	Script  string `json:"script"`
	Addr    string `json:"addr"`
}

type UnspentOutput struct {
	TxAge           int    `json:"tx_age"`
	TxHash          string `json:"tx_hash"`
	TxHashBigEndian string `json:"tx_hash_big_endian"`
	TxIndex         int64  `json:"tx_index"`
	TxOutputN       int    `json:"tx_output_n"`
	Script          string `json:"script"`
	Value           int64  `json:"value"`
	ValueHex        string `json:"value_hex"`
	Confirmations   int    `json:"confirmations"`
}

type UnspentOutputs struct {
	UnspentOutputs []UnspentOutput `json:"unspent_outputs"`
}

type Balance struct {
	Address       string `json:"address"`
	Balance       int64  `json:"balance"`
	TotalReceived int64  `json:"total_received"`
	TotalSent     int64  `json:"total_sent"`
}

type MultiAddress struct {
	Addresses []AddressInfo `json:"addresses"`
	Wallet    struct {
		NTx           int   `json:"n_tx"`
		NTxFiltered   int   `json:"n_tx_filtered"`
		TotalReceived int64 `json:"total_received"`
		TotalSent     int64 `json:"total_sent"`
		FinalBalance  int64 `json:"final_balance"`
	} `json:"wallet"`
	Txs  []Tx `json:"txs"`
	Info struct {
		NConnected int     `json:"nconnected"`
		Conversion float64 `json:"conversion"`
		Symbol     struct {
			Code string `json:"code"`
		} `json:"symbol_local"`
		SymbolBTC struct {
			Code string `json:"code"`
		} `json:"symbol_btc"`
		Latest struct {
			BlockIndex int64  `json:"block_index"`
			Hash       string `json:"hash"`
			Height     int    `json:"height"`
			Time       int64  `json:"time"`
		} `json:"latest_block"`
	} `json:"info"`
}

type PushTxRequest struct {
	Tx string `json:"tx"`
}

type PushTxResponse struct {
	Notice string `json:"notice,omitempty"`
	Error  string `json:"error,omitempty"`
}
