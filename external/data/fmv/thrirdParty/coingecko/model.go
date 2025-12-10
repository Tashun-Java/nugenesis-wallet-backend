package coingecko

// PriceResponse represents the response from CoinGecko simple/price endpoint
type PriceResponse map[string]map[string]float64

// symbolToCoinGeckoID maps common token symbols to their CoinGecko IDs
var SymbolToCoinGeckoID = map[string]string{
	"BTC":    "bitcoin",
	"ETH":    "ethereum",
	"USDT":   "tether",
	"USDC":   "usd-coin",
	"BNB":    "binancecoin",
	"SOL":    "solana",
	"MATIC":  "matic-network",
	"POL":    "matic-network",
	"AVAX":   "avalanche-2",
	"DAI":    "dai",
	"WBTC":   "wrapped-bitcoin",
	"WETH":   "weth",
	"LINK":   "chainlink",
	"UNI":    "uniswap",
	"AAVE":   "aave",
	"USDS":   "usds",
	"MNDE":   "marinade",
	"MSOL":   "msol",
	"BUSD":   "binance-usd",
	"XRP":    "ripple",
	"ADA":    "cardano",
	"DOGE":   "dogecoin",
	"TRX":    "tron",
	"DOT":    "polkadot",
	"SHIB":   "shiba-inu",
	"LTC":    "litecoin",
	"BCH":    "bitcoin-cash",
	"ATOM":   "cosmos",
	"ARB":    "arbitrum",
	"OP":     "optimism",
	"BASE":   "base",
	"BLAST":  "blast",
	"ZKSYNC": "zksync",
	"SCROLL": "scroll",
}
