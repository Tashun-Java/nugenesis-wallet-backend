package models

type Transaction struct {
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

//type BlockchainDecimal string
//
//const (
//	Bitcoin          BlockchainDecimal = "0"
//	Ethereum         BlockchainDecimal = "1"
//	Vechain          BlockchainDecimal = "2"
//	Tron             BlockchainDecimal = "3"
//	Icon             BlockchainDecimal = "4"
//	Binance          BlockchainDecimal = "5"
//	Ripple           BlockchainDecimal = "6"
//	Tezos            BlockchainDecimal = "7"
//	Nimiq            BlockchainDecimal = "8"
//	Stellar          BlockchainDecimal = "9"
//	Aion             BlockchainDecimal = "10"
//	Cosmos           BlockchainDecimal = "11"
//	Theta            BlockchainDecimal = "12"
//	Ontology         BlockchainDecimal = "13"
//	Zilliqa          BlockchainDecimal = "14"
//	IoTeX            BlockchainDecimal = "15"
//	Eos              BlockchainDecimal = "16"
//	Nano             BlockchainDecimal = "17"
//	Nuls             BlockchainDecimal = "18"
//	Waves            BlockchainDecimal = "19"
//	Aeternity        BlockchainDecimal = "20"
//	Nebulas          BlockchainDecimal = "21"
//	Fio              BlockchainDecimal = "22"
//	Solana           BlockchainDecimal = "23"
//	Harmony          BlockchainDecimal = "24"
//	Near             BlockchainDecimal = "25"
//	Algorand         BlockchainDecimal = "26"
//	Iost             BlockchainDecimal = "27"
//	Polkadot         BlockchainDecimal = "28"
//	Cardano          BlockchainDecimal = "29"
//	Neo              BlockchainDecimal = "30"
//	Filecoin         BlockchainDecimal = "31"
//	MultiversX       BlockchainDecimal = "32"
//	OasisNetwork     BlockchainDecimal = "33"
//	Decred           BlockchainDecimal = "34"
//	Zcash            BlockchainDecimal = "35"
//	Groestlcoin      BlockchainDecimal = "36"
//	Thorchain        BlockchainDecimal = "37"
//	Ronin            BlockchainDecimal = "38"
//	Kusama           BlockchainDecimal = "39"
//	Zen              BlockchainDecimal = "40"
//	BitcoinDiamond   BlockchainDecimal = "41"
//	Verge            BlockchainDecimal = "42"
//	Nervos           BlockchainDecimal = "43"
//	Everscale        BlockchainDecimal = "44"
//	Aptos            BlockchainDecimal = "45"
//	Nebl             BlockchainDecimal = "46"
//	Hedera           BlockchainDecimal = "47"
//	TheOpenNetwork   BlockchainDecimal = "48"
//	Sui              BlockchainDecimal = "49"
//	Greenfield       BlockchainDecimal = "50"
//	InternetComputer BlockchainDecimal = "51"
//	NativeEvmos      BlockchainDecimal = "52"
//	NativeInjective  BlockchainDecimal = "53"
//	BitcoinCash      BlockchainDecimal = "54"
//	Pactus           BlockchainDecimal = "55"
//	Komodo           BlockchainDecimal = "56"
//	Polymesh         BlockchainDecimal = "57"
//	Sepolia          BlockchainDecimal = "58"
//)
