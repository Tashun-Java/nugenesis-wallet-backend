package staticServices

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tashunc/nugenesis-wallet-backend/external/models/general"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticModels"
)

// BlockchainService handles blockchain ID mapping and operations
type BlockchainService struct {
	blockchainMap map[string]general.CoinType // blockchain name -> CoinType ID
	assetService  *AssetService
	mutex         sync.RWMutex
}

// NewBlockchainService creates a new BlockchainService instance
func NewBlockchainService(assetService *AssetService) *BlockchainService {
	service := &BlockchainService{
		blockchainMap: make(map[string]general.CoinType),
		assetService:  assetService,
	}
	service.initializeBlockchainMapping()
	return service
}

// initializeBlockchainMapping creates the mapping between blockchain names and CoinType IDs
func (s *BlockchainService) initializeBlockchainMapping() {
	// Main networks
	s.blockchainMap["aeternity"] = general.Aeternity
	s.blockchainMap["aion"] = general.Aion
	s.blockchainMap["binance"] = general.Binance
	s.blockchainMap["bitcoin"] = general.Bitcoin
	s.blockchainMap["bitcoincash"] = general.BitcoinCash
	s.blockchainMap["bitcoingold"] = general.BitcoinGold
	s.blockchainMap["callisto"] = general.Callisto
	s.blockchainMap["cardano"] = general.Cardano
	s.blockchainMap["cosmos"] = general.Cosmos
	s.blockchainMap["pivx"] = general.Pivx
	s.blockchainMap["dash"] = general.Dash
	s.blockchainMap["decred"] = general.Decred
	s.blockchainMap["digibyte"] = general.DigiByte
	s.blockchainMap["dogecoin"] = general.Dogecoin
	s.blockchainMap["eos"] = general.Eos
	s.blockchainMap["wax"] = general.Wax
	s.blockchainMap["ethereum"] = general.Ethereum
	s.blockchainMap["ethereumclassic"] = general.EthereumClassic
	s.blockchainMap["fio"] = general.Fio
	s.blockchainMap["gochain"] = general.GoChain
	s.blockchainMap["groestlcoin"] = general.Groestlcoin
	s.blockchainMap["icon"] = general.Icon
	s.blockchainMap["iotex"] = general.IoTeX
	s.blockchainMap["kava"] = general.Kava
	s.blockchainMap["kin"] = general.Kin
	s.blockchainMap["litecoin"] = general.Litecoin
	s.blockchainMap["monacoin"] = general.Monacoin
	s.blockchainMap["nebulas"] = general.Nebulas
	s.blockchainMap["nuls"] = general.Nuls
	s.blockchainMap["nano"] = general.Nano
	s.blockchainMap["near"] = general.Near
	s.blockchainMap["nimiq"] = general.Nimiq
	s.blockchainMap["ontology"] = general.Ontology
	s.blockchainMap["poanetwork"] = general.Poanetwork
	s.blockchainMap["qtum"] = general.Qtum
	s.blockchainMap["xrp"] = general.Xrp
	s.blockchainMap["solana"] = general.Solana
	s.blockchainMap["stellar"] = general.Stellar
	s.blockchainMap["tezos"] = general.Tezos
	s.blockchainMap["theta"] = general.Theta
	s.blockchainMap["thundercore"] = general.ThunderCore
	s.blockchainMap["neo"] = general.Neo
	s.blockchainMap["viction"] = general.Viction
	s.blockchainMap["tron"] = general.Tron
	s.blockchainMap["vechain"] = general.VeChain
	s.blockchainMap["viacoin"] = general.Viacoin
	s.blockchainMap["wanchain"] = general.Wanchain
	s.blockchainMap["zcash"] = general.Zcash
	s.blockchainMap["firo"] = general.Firo
	s.blockchainMap["zilliqa"] = general.Zilliqa
	s.blockchainMap["zelcash"] = general.Zelcash
	s.blockchainMap["ravencoin"] = general.Ravencoin
	s.blockchainMap["waves"] = general.Waves
	s.blockchainMap["terra"] = general.Terra
	s.blockchainMap["terrav2"] = general.TerraV2
	s.blockchainMap["harmony"] = general.Harmony
	s.blockchainMap["algorand"] = general.Algorand
	s.blockchainMap["kusama"] = general.Kusama
	s.blockchainMap["polkadot"] = general.Polkadot
	s.blockchainMap["filecoin"] = general.Filecoin
	s.blockchainMap["multiversx"] = general.MultiversX
	s.blockchainMap["bandchain"] = general.BandChain
	s.blockchainMap["smartchainlegacy"] = general.SmartChainLegacy
	s.blockchainMap["smartchain"] = general.SmartChain
	s.blockchainMap["tbinance"] = general.TBinance
	s.blockchainMap["oasis"] = general.Oasis
	s.blockchainMap["polygon"] = general.Polygon
	s.blockchainMap["thorchain"] = general.Thorchain
	s.blockchainMap["bluzelle"] = general.Bluzelle
	s.blockchainMap["optimism"] = general.Optimism
	s.blockchainMap["zksync"] = general.Zksync
	s.blockchainMap["arbitrum"] = general.Arbitrum
	s.blockchainMap["ecochain"] = general.Ecochain
	s.blockchainMap["avalanchec"] = general.AvalancheCChain
	s.blockchainMap["xdai"] = general.Xdai
	s.blockchainMap["fantom"] = general.Fantom
	s.blockchainMap["cryptoorg"] = general.CryptoOrg
	s.blockchainMap["celo"] = general.Celo
	s.blockchainMap["ronin"] = general.Ronin
	s.blockchainMap["osmosis"] = general.Osmosis
	s.blockchainMap["ecash"] = general.Ecash
	s.blockchainMap["iost"] = general.Iost
	s.blockchainMap["cronos"] = general.CronosChain
	s.blockchainMap["smartbch"] = general.SmartBitcoinCash
	s.blockchainMap["kcc"] = general.KuCoinCommunityChain
	s.blockchainMap["bitcoindiamond"] = general.BitcoinDiamond
	s.blockchainMap["boba"] = general.Boba
	s.blockchainMap["syscoin"] = general.Syscoin
	s.blockchainMap["verge"] = general.Verge
	s.blockchainMap["zen"] = general.Zen
	s.blockchainMap["metis"] = general.Metis
	s.blockchainMap["aurora"] = general.Aurora
	s.blockchainMap["evmos"] = general.Evmos
	s.blockchainMap["nativeevmos"] = general.NativeEvmos
	s.blockchainMap["moonriver"] = general.Moonriver
	s.blockchainMap["moonbeam"] = general.Moonbeam
	s.blockchainMap["kavaevm"] = general.KavaEvm
	s.blockchainMap["kaia"] = general.Kaia
	s.blockchainMap["meter"] = general.Meter
	s.blockchainMap["okxchain"] = general.Okxchain
	s.blockchainMap["stratis"] = general.Stratis
	s.blockchainMap["komodo"] = general.Komodo
	s.blockchainMap["nervos"] = general.Nervos
	s.blockchainMap["everscale"] = general.Everscale
	s.blockchainMap["aptos"] = general.Aptos
	s.blockchainMap["nebl"] = general.Nebl
	s.blockchainMap["hedera"] = general.Hedera
	s.blockchainMap["secret"] = general.Secret
	s.blockchainMap["nativeinjective"] = general.NativeInjective
	s.blockchainMap["agoric"] = general.Agoric
	s.blockchainMap["ton"] = general.Ton
	s.blockchainMap["sui"] = general.Sui
	s.blockchainMap["stargaze"] = general.Stargaze
	s.blockchainMap["polygonzkevm"] = general.PolygonzkEVM
	s.blockchainMap["juno"] = general.Juno
	s.blockchainMap["stride"] = general.Stride
	s.blockchainMap["axelar"] = general.Axelar
	s.blockchainMap["crescent"] = general.Crescent
	s.blockchainMap["kujira"] = general.Kujira
	s.blockchainMap["iotexevm"] = general.IoTeXEVM
	s.blockchainMap["nativecanto"] = general.NativeCanto
	s.blockchainMap["comdex"] = general.Comdex
	s.blockchainMap["neutron"] = general.Neutron
	s.blockchainMap["sommelier"] = general.Sommelier
	s.blockchainMap["fetchai"] = general.FetchAI
	s.blockchainMap["mars"] = general.Mars
	s.blockchainMap["umee"] = general.Umee
	s.blockchainMap["coreum"] = general.Coreum
	s.blockchainMap["quasar"] = general.Quasar
	s.blockchainMap["persistence"] = general.Persistence
	s.blockchainMap["akash"] = general.Akash
	s.blockchainMap["noble"] = general.Noble
	s.blockchainMap["scroll"] = general.Scroll
	s.blockchainMap["rootstock"] = general.Rootstock
	s.blockchainMap["thetafuel"] = general.ThetaFuel
	s.blockchainMap["cfxevm"] = general.ConfluxeSpace
	s.blockchainMap["acala"] = general.Acala
	s.blockchainMap["acalaevm"] = general.AcalaEVM
	s.blockchainMap["opbnb"] = general.OpBNB
	s.blockchainMap["neon"] = general.Neon
	s.blockchainMap["base"] = general.Base
	s.blockchainMap["sei"] = general.Sei
	s.blockchainMap["arbitrumnova"] = general.ArbitrumNova
	s.blockchainMap["linea"] = general.Linea
	s.blockchainMap["greenfield"] = general.Greenfield
	s.blockchainMap["mantle"] = general.Mantle
	s.blockchainMap["zeneon"] = general.ZenEON
	s.blockchainMap["internetcomputer"] = general.InternetComputer
	s.blockchainMap["tia"] = general.Tia
	s.blockchainMap["manta"] = general.MantaPacific
	s.blockchainMap["nativezetachain"] = general.NativeZetaChain
	s.blockchainMap["zetaevm"] = general.ZetaEVM
	s.blockchainMap["dydx"] = general.Dydx
	s.blockchainMap["merlin"] = general.Merlin
	s.blockchainMap["lightlink"] = general.Lightlink
	s.blockchainMap["blast"] = general.Blast
	s.blockchainMap["bouncebit"] = general.BounceBit
	s.blockchainMap["zklink"] = general.ZkLinkNova
	s.blockchainMap["pactus"] = general.Pactus
	s.blockchainMap["sonic"] = general.Sonic
	s.blockchainMap["polymesh"] = general.Polymesh

}

// GetBlockchainID returns the CoinType ID for a given blockchain name
func (s *BlockchainService) GetBlockchainID(blockchainName string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if coinType, exists := s.blockchainMap[strings.ToLower(blockchainName)]; exists {
		return string(coinType), true
	}
	return "", false
}

// GetAllBlockchains returns all blockchains with their IDs and asset counts
func (s *BlockchainService) GetAllBlockchains() ([]staticModels.BlockchainResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Get available blockchains from assets directory
	availableBlockchains := s.getAvailableBlockchains()

	var blockchains []staticModels.BlockchainResponse

	for blockchainName, coinType := range s.blockchainMap {
		// Check if this blockchain has assets
		hasAssets := false
		assetCount := 0

		if info, exists := availableBlockchains[blockchainName]; exists {
			hasAssets = true
			assetCount = info.AssetCount
		}

		// Determine if it's a testnet
		isTestnet := strings.Contains(strings.ToLower(blockchainName), "test") ||
			strings.Contains(strings.ToLower(blockchainName), "goerli") ||
			strings.Contains(strings.ToLower(blockchainName), "mumbai") ||
			strings.Contains(strings.ToLower(blockchainName), "sepolia")

		blockchain := staticModels.BlockchainResponse{
			ID:         string(coinType),
			Name:       blockchainName,
			IsTestnet:  isTestnet,
			HasAssets:  hasAssets,
			AssetCount: assetCount,
		}

		// Read blockchain info.json to get additional details
		infoPath := filepath.Join("assets/blockchains", blockchainName, "info", "info.json")
		if infoData, err := ioutil.ReadFile(infoPath); err == nil {
			var blockchainInfo staticModels.BlockchainInfo
			if err := json.Unmarshal(infoData, &blockchainInfo); err == nil {
				blockchain.Symbol = blockchainInfo.Symbol
				blockchain.Website = blockchainInfo.Website
				blockchain.Explorer = blockchainInfo.Explorer
				blockchain.Decimals = blockchainInfo.Decimals
				blockchain.ImageURL = "/static/assetsLogo/blockchains/" + blockchainName + "/info/logo.png"
			}
		}

		blockchains = append(blockchains, blockchain)
	}

	return blockchains, nil
}

// GetBlockchainByID returns a blockchain by its CoinType ID
func (s *BlockchainService) GetBlockchainByID(id string) (*staticModels.BlockchainResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Find blockchain by ID
	for blockchainName, coinType := range s.blockchainMap {
		if string(coinType) == id {
			// Get asset count for this blockchain
			availableBlockchains := s.getAvailableBlockchains()
			hasAssets := false
			assetCount := 0

			if info, exists := availableBlockchains[blockchainName]; exists {
				hasAssets = true
				assetCount = info.AssetCount
			}

			// Determine if it's a testnet
			isTestnet := strings.Contains(strings.ToLower(blockchainName), "test") ||
				strings.Contains(strings.ToLower(blockchainName), "goerli") ||
				strings.Contains(strings.ToLower(blockchainName), "mumbai") ||
				strings.Contains(strings.ToLower(blockchainName), "sepolia")

			blockchain := &staticModels.BlockchainResponse{
				ID:         id,
				Name:       blockchainName,
				IsTestnet:  isTestnet,
				HasAssets:  hasAssets,
				AssetCount: assetCount,
			}

			// Read blockchain info.json to get additional details
			infoPath := filepath.Join("assets/blockchains", blockchainName, "info", "info.json")
			if infoData, err := ioutil.ReadFile(infoPath); err == nil {
				var blockchainInfo staticModels.BlockchainInfo
				if err := json.Unmarshal(infoData, &blockchainInfo); err == nil {
					blockchain.Symbol = blockchainInfo.Symbol
					blockchain.Website = blockchainInfo.Website
					blockchain.Explorer = blockchainInfo.Explorer
					blockchain.Decimals = blockchainInfo.Decimals
					blockchain.ImageURL = "/static/assetsLogo/blockchains/" + blockchainName + "/info/logo.png"
				}
			}

			return blockchain, nil
		}
	}

	return nil, nil // Not found
}

// getAvailableBlockchains scans the assets directory to find available blockchains
func (s *BlockchainService) getAvailableBlockchains() map[string]struct {
	AssetCount int
} {
	result := make(map[string]struct {
		AssetCount int
	})

	// Read blockchain directories
	entries, err := ioutil.ReadDir("assets/blockchains")
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if entry.IsDir() {
			blockchainName := entry.Name()
			assetCount := 0

			// Count assets for this blockchain
			// Check if native token exists
			infoPath := filepath.Join("assets/blockchains", blockchainName, "info", "info.json")
			if _, err := ioutil.ReadFile(infoPath); err == nil {
				assetCount++
			}

			// Count token assets
			assetsDir := filepath.Join("assets/blockchains", blockchainName, "assets")
			if assetEntries, err := ioutil.ReadDir(assetsDir); err == nil {
				for _, assetEntry := range assetEntries {
					if assetEntry.IsDir() {
						assetInfoPath := filepath.Join(assetsDir, assetEntry.Name(), "info.json")
						if _, err := ioutil.ReadFile(assetInfoPath); err == nil {
							assetCount++
						}
					}
				}
			}

			result[blockchainName] = struct {
				AssetCount int
			}{AssetCount: assetCount}
		}
	}

	return result
}
