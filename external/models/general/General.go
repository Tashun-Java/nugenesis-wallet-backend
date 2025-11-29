package general

type CoinType string

const (
	Aeternity                CoinType = "457"
	Aion                     CoinType = "425"
	Binance                  CoinType = "714"
	Bitcoin                  CoinType = "0"
	BitcoinCash              CoinType = "145"
	BitcoinGold              CoinType = "156"
	Callisto                 CoinType = "820"
	Cardano                  CoinType = "1815"
	Cosmos                   CoinType = "118"
	Pivx                     CoinType = "119"
	Dash                     CoinType = "5"
	Decred                   CoinType = "42"
	DigiByte                 CoinType = "20"
	Dogecoin                 CoinType = "3"
	Eos                      CoinType = "194"
	Wax                      CoinType = "14001"
	Ethereum                 CoinType = "60"
	EthereumClassic          CoinType = "61"
	Fio                      CoinType = "235"
	GoChain                  CoinType = "6060"
	Groestlcoin              CoinType = "17"
	Icon                     CoinType = "74"
	IoTeX                    CoinType = "304"
	Kava                     CoinType = "459"
	Kin                      CoinType = "2017"
	Litecoin                 CoinType = "2"
	Monacoin                 CoinType = "22"
	Nebulas                  CoinType = "2718"
	Nuls                     CoinType = "8964"
	Nano                     CoinType = "165"
	Near                     CoinType = "397"
	Nimiq                    CoinType = "242"
	Ontology                 CoinType = "1024"
	Poanetwork               CoinType = "178"
	Qtum                     CoinType = "2301"
	Xrp                      CoinType = "144"
	Solana                   CoinType = "501"
	Stellar                  CoinType = "148"
	Tezos                    CoinType = "1729"
	Theta                    CoinType = "500"
	ThunderCore              CoinType = "1001"
	Neo                      CoinType = "888"
	Viction                  CoinType = "889"
	Tron                     CoinType = "195"
	VeChain                  CoinType = "818"
	Viacoin                  CoinType = "14"
	Wanchain                 CoinType = "5718350"
	Zcash                    CoinType = "133"
	Firo                     CoinType = "136"
	Zilliqa                  CoinType = "313"
	Zelcash                  CoinType = "19167"
	Ravencoin                CoinType = "175"
	Waves                    CoinType = "5741564"
	Terra                    CoinType = "330"
	TerraV2                  CoinType = "10000330"
	Harmony                  CoinType = "1023"
	Algorand                 CoinType = "283"
	Kusama                   CoinType = "434"
	Polkadot                 CoinType = "354"
	Filecoin                 CoinType = "461"
	MultiversX               CoinType = "508"
	BandChain                CoinType = "494"
	SmartChainLegacy         CoinType = "10000714"
	SmartChain               CoinType = "20000714"
	TBinance                 CoinType = "30000714"
	Oasis                    CoinType = "474"
	Polygon                  CoinType = "966"
	Thorchain                CoinType = "931"
	Bluzelle                 CoinType = "483"
	Optimism                 CoinType = "10000070"
	Zksync                   CoinType = "10000324"
	Arbitrum                 CoinType = "10042221"
	Ecochain                 CoinType = "10000553"
	AvalancheCChain          CoinType = "10009000"
	Xdai                     CoinType = "10000100"
	Fantom                   CoinType = "10000250"
	CryptoOrg                CoinType = "394"
	Celo                     CoinType = "52752"
	Ronin                    CoinType = "10002020"
	Osmosis                  CoinType = "10000118"
	Ecash                    CoinType = "899"
	Iost                     CoinType = "291"
	CronosChain              CoinType = "10000025"
	SmartBitcoinCash         CoinType = "10000145"
	KuCoinCommunityChain     CoinType = "10000321"
	BitcoinDiamond           CoinType = "999"
	Boba                     CoinType = "10000288"
	Syscoin                  CoinType = "57"
	Verge                    CoinType = "77"
	Zen                      CoinType = "121"
	Metis                    CoinType = "10001088"
	Aurora                   CoinType = "1323161554"
	Evmos                    CoinType = "10009001"
	NativeEvmos              CoinType = "20009001"
	Moonriver                CoinType = "10001285"
	Moonbeam                 CoinType = "10001284"
	KavaEvm                  CoinType = "10002222"
	Kaia                     CoinType = "10008217"
	Meter                    CoinType = "18000"
	Okxchain                 CoinType = "996"
	Stratis                  CoinType = "105105"
	Komodo                   CoinType = "141"
	Nervos                   CoinType = "309"
	Everscale                CoinType = "396"
	Aptos                    CoinType = "637"
	Nebl                     CoinType = "146"
	Hedera                   CoinType = "3030"
	Secret                   CoinType = "529"
	NativeInjective          CoinType = "10000060"
	Agoric                   CoinType = "564"
	Ton                      CoinType = "607"
	Sui                      CoinType = "784"
	Stargaze                 CoinType = "20000118"
	PolygonzkEVM             CoinType = "10001101"
	Juno                     CoinType = "30000118"
	Stride                   CoinType = "40000118"
	Axelar                   CoinType = "50000118"
	Crescent                 CoinType = "60000118"
	Kujira                   CoinType = "70000118"
	IoTeXEVM                 CoinType = "10004689"
	NativeCanto              CoinType = "10007700"
	Comdex                   CoinType = "80000118"
	Neutron                  CoinType = "90000118"
	Sommelier                CoinType = "11000118"
	FetchAI                  CoinType = "12000118"
	Mars                     CoinType = "13000118"
	Umee                     CoinType = "14000118"
	Coreum                   CoinType = "10000990"
	Quasar                   CoinType = "15000118"
	Persistence              CoinType = "16000118"
	Akash                    CoinType = "17000118"
	Noble                    CoinType = "18000118"
	Scroll                   CoinType = "534352"
	Rootstock                CoinType = "137"
	ThetaFuel                CoinType = "361"
	ConfluxeSpace            CoinType = "1030"
	Acala                    CoinType = "787"
	AcalaEVM                 CoinType = "10000787"
	OpBNB                    CoinType = "204"
	Neon                     CoinType = "245022934"
	Base                     CoinType = "8453"
	Sei                      CoinType = "19000118"
	ArbitrumNova             CoinType = "10042170"
	Linea                    CoinType = "59144"
	Greenfield               CoinType = "5600"
	Mantle                   CoinType = "5000"
	ZenEON                   CoinType = "7332"
	InternetComputer         CoinType = "223"
	Tia                      CoinType = "21000118"
	MantaPacific             CoinType = "169"
	NativeZetaChain          CoinType = "10007000"
	ZetaEVM                  CoinType = "20007000"
	Dydx                     CoinType = "22000118"
	Merlin                   CoinType = "4200"
	Lightlink                CoinType = "1890"
	Blast                    CoinType = "81457"
	BounceBit                CoinType = "6001"
	ZkLinkNova               CoinType = "810180"
	Pactus                   CoinType = "21888"
	Sonic                    CoinType = "10000146"
	Polymesh                 CoinType = "595"

)

//export class EthereumChainID {
//value: number;
//static ethereum: EthereumChainID;
//static classic: EthereumChainID;
//static rootstock: EthereumChainID;
//static manta: EthereumChainID;
//static poa: EthereumChainID;
//static opbnb: EthereumChainID;
//static tfuelevm: EthereumChainID;
//static vechain: EthereumChainID;
//static callisto: EthereumChainID;
//static viction: EthereumChainID;
//static polygon: EthereumChainID;
//static okc: EthereumChainID;
//static thundertoken: EthereumChainID;
//static cfxevm: EthereumChainID;
//static lightlink: EthereumChainID;
//static merlin: EthereumChainID;
//static mantle: EthereumChainID;
//static bouncebit: EthereumChainID;
//static gochain: EthereumChainID;
//static zeneon: EthereumChainID;
//static base: EthereumChainID;
//static meter: EthereumChainID;
//static celo: EthereumChainID;
//static linea: EthereumChainID;
//static blast: EthereumChainID;
//static scroll: EthereumChainID;
//static zklinknova: EthereumChainID;
//static wanchain: EthereumChainID;
//static cronos: EthereumChainID;
//static optimism: EthereumChainID;
//static xdai: EthereumChainID;
//static smartbch: EthereumChainID;
//static sonic: EthereumChainID;
//static fantom: EthereumChainID;
//static boba: EthereumChainID;
//static kcc: EthereumChainID;
//static zksync: EthereumChainID;
//static heco: EthereumChainID;
//static acalaevm: EthereumChainID;
//static metis: EthereumChainID;
//static polygonzkevm: EthereumChainID;
//static moonbeam: EthereumChainID;
//static moonriver: EthereumChainID;
//static ronin: EthereumChainID;
//static kavaevm: EthereumChainID;
//static iotexevm: EthereumChainID;
//static kaia: EthereumChainID;
//static avalanchec: EthereumChainID;
//static evmos: EthereumChainID;
//static arbitrumnova: EthereumChainID;
//static arbitrum: EthereumChainID;
//static smartchain: EthereumChainID;
//static zetaevm: EthereumChainID;
//static neon: EthereumChainID;
//static aurora: EthereumChainID;
//}
