package staticServices

import (
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/models/general"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticModels"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AssetIDMapping represents the mapping between asset keys and their persistent IDs
type AssetIDMapping struct {
	AssetKey string `json:"asset_key"`
	ID       string `json:"id"`
}

// AssetService handles blockchain assets management
type AssetService struct {
	assetCache       map[string][]staticModels.AssetResponse
	idMappings       map[string]string           // asset_key -> id
	solanaTokenCache map[string]string           // Solana mint address -> token symbol
	nextID           int                         // counter for generating new IDs
	blockchainMap    map[string]general.CoinType // blockchain name -> CoinType ID
	mutex            sync.RWMutex
	tokenCacheMutex  sync.RWMutex
	lastUpdate       time.Time
	cacheTTL         time.Duration
	assetsPath       string
	idMappingFile    string
	tokenCacheLoaded bool
}

// NewAssetService creates a new AssetService instance
func NewAssetService() *AssetService {
	service := &AssetService{
		assetCache:       make(map[string][]staticModels.AssetResponse),
		idMappings:       make(map[string]string),
		solanaTokenCache: make(map[string]string),
		blockchainMap:    make(map[string]general.CoinType),
		nextID:           1,
		cacheTTL:         30 * time.Minute, // Cache for 30 minutes
		assetsPath:       "./assets/blockchains",
		idMappingFile:    "./assets/id_mappings.json",
		tokenCacheLoaded: false,
	}
	service.initializeBlockchainMapping()
	service.loadIDMappings()
	return service
}

// initializeBlockchainMapping creates the mapping between blockchain names and CoinType IDs
func (s *AssetService) initializeBlockchainMapping() {
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
	// Add more mappings as needed
}

// getBlockchainID returns the CoinType ID for a given blockchain name
func (s *AssetService) getBlockchainID(blockchainName string) string {
	if coinType, exists := s.blockchainMap[strings.ToLower(blockchainName)]; exists {
		return string(coinType)
	}
	return ""
}

// generateAssetKey creates a unique key for an asset
func (s *AssetService) generateAssetKey(blockchain, address, symbol string) string {
	if address != "" {
		return fmt.Sprintf("%s-%s", blockchain, address)
	}
	return fmt.Sprintf("%s-%s-native", blockchain, symbol)
}

// getOrCreateAssetID gets existing ID or creates new one for an asset
func (s *AssetService) getOrCreateAssetID(blockchain, address, symbol string) string {
	assetKey := s.generateAssetKey(blockchain, address, symbol)

	// Check if ID already exists
	if id, exists := s.idMappings[assetKey]; exists {
		return id
	}

	// Create new ID
	newID := strconv.Itoa(s.nextID)
	s.idMappings[assetKey] = newID
	s.nextID++

	// Don't save on every asset - we'll save in batch later
	return newID
}

// loadIDMappings loads ID mappings from file
func (s *AssetService) loadIDMappings() {
	data, err := ioutil.ReadFile(s.idMappingFile)
	if err != nil {
		// File doesn't exist yet, start with empty mappings
		return
	}

	var mappings []AssetIDMapping
	if err := json.Unmarshal(data, &mappings); err != nil {
		fmt.Printf("Error loading ID mappings: %v\n", err)
		return
	}

	// Convert to map and find next ID
	maxID := 0
	for _, mapping := range mappings {
		s.idMappings[mapping.AssetKey] = mapping.ID
		if id, err := strconv.Atoi(mapping.ID); err == nil && id > maxID {
			maxID = id
		}
	}
	s.nextID = maxID + 1
}

// saveIDMappings saves ID mappings to file
func (s *AssetService) saveIDMappings() {
	if len(s.idMappings) == 0 {
		return // Nothing to save
	}

	var mappings []AssetIDMapping
	for key, id := range s.idMappings {
		mappings = append(mappings, AssetIDMapping{
			AssetKey: key,
			ID:       id,
		})
	}

	data, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling ID mappings: %v\n", err)
		return
	}

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(s.idMappingFile), 0755); err != nil {
		fmt.Printf("Error creating directory for ID mappings: %v\n", err)
		return
	}

	if err := ioutil.WriteFile(s.idMappingFile, data, 0644); err != nil {
		fmt.Printf("Error saving ID mappings to %s: %v\n", s.idMappingFile, err)
	}
}

// loadBlockchainAssets loads all assets for a specific blockchain
func (s *AssetService) loadBlockchainAssets(blockchainPath, blockchainName string) ([]staticModels.AssetResponse, error) {
	var assets []staticModels.AssetResponse

	// Load native token info
	infoPath := filepath.Join(blockchainPath, "info", "info.json")
	logoPath := filepath.Join(blockchainPath, "info", "logo.png")

	if infoData, err := ioutil.ReadFile(infoPath); err == nil {
		var blockchainInfo staticModels.BlockchainInfo
		if err := json.Unmarshal(infoData, &blockchainInfo); err == nil {
			// Check if logo exists
			logoExists := false
			if _, err := os.Stat(logoPath); err == nil {
				logoExists = true
			}

			// Skip if blockchain ID mapping doesn't exist
			blockchainID := s.getBlockchainID(blockchainName)
			if blockchainID == "" {
				return assets, nil
			}

			nativeAsset := staticModels.AssetResponse{
				ID:           s.getOrCreateAssetID(blockchainName, "", blockchainInfo.Symbol),
				Symbol:       blockchainInfo.Symbol,
				Name:         blockchainInfo.Name,
				Blockchain:   blockchainName,
				BlockchainID: blockchainID,
				Address:      "", // Native tokens don't have contract addresses
				Type:         blockchainInfo.Type,
				Decimals:     blockchainInfo.Decimals,
				Status:       blockchainInfo.Status,
			}

			if logoExists {
				nativeAsset.LogoPath = fmt.Sprintf("/static/assetsLogo/blockchains/%s/info/logo.png", blockchainName)
			}

			assets = append(assets, nativeAsset)
		}
	}

	// Load token assets
	assetsDir := filepath.Join(blockchainPath, "assets")
	if entries, err := ioutil.ReadDir(assetsDir); err == nil {
		// Get blockchain ID once to check if we should process tokens
		blockchainID := s.getBlockchainID(blockchainName)
		if blockchainID == "" {
			return assets, nil
		}

		for _, entry := range entries {
			if entry.IsDir() {
				assetPath := filepath.Join(assetsDir, entry.Name())
				assetInfoPath := filepath.Join(assetPath, "info.json")
				assetLogoPath := filepath.Join(assetPath, "logo.png")

				if assetData, err := ioutil.ReadFile(assetInfoPath); err == nil {
					var assetInfo staticModels.AssetInfo
					if err := json.Unmarshal(assetData, &assetInfo); err == nil {
						// Check if logo exists
						logoExists := false
						if _, err := os.Stat(assetLogoPath); err == nil {
							logoExists = true
						}

						tokenAsset := staticModels.AssetResponse{
							ID:           s.getOrCreateAssetID(blockchainName, entry.Name(), assetInfo.Symbol),
							Symbol:       assetInfo.Symbol,
							Name:         assetInfo.Name,
							Blockchain:   blockchainName,
							BlockchainID: blockchainID,
							Address:      entry.Name(), // Directory name is the contract address/asset ID
							Type:         assetInfo.Type,
							Decimals:     assetInfo.Decimals,
							Status:       assetInfo.Status,
						}

						if logoExists {
							tokenAsset.LogoPath = fmt.Sprintf("/static/assetsLogo/blockchains/%s/assets/%s/logo.png", blockchainName, entry.Name())
						}

						assets = append(assets, tokenAsset)
					}
				}
			}
		}
	}

	return assets, nil
}

// loadAllAssets loads all blockchain assets and populates the cache
func (s *AssetService) loadAllAssets() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Clear existing cache
	s.assetCache = make(map[string][]staticModels.AssetResponse)
	initialMappingCount := len(s.idMappings)

	// Read all blockchain directories
	entries, err := ioutil.ReadDir(s.assetsPath)
	if err != nil {
		return fmt.Errorf("failed to read assets directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			blockchainName := entry.Name()

			// Skip if blockchain ID mapping doesn't exist
			if s.getBlockchainID(blockchainName) == "" {
				continue
			}

			blockchainPath := filepath.Join(s.assetsPath, blockchainName)

			assets, err := s.loadBlockchainAssets(blockchainPath, blockchainName)
			if err != nil {
				continue // Skip this blockchain on error but continue with others
			}

			// Group assets by symbol
			for _, asset := range assets {
				symbol := strings.ToUpper(asset.Symbol)
				s.assetCache[symbol] = append(s.assetCache[symbol], asset)
			}
		}
	}

	// Save ID mappings only if new ones were created
	if len(s.idMappings) > initialMappingCount {
		s.saveIDMappings()
	}

	s.lastUpdate = time.Now()
	return nil
}

// refreshCacheIfNeeded refreshes the cache if it's expired
func (s *AssetService) refreshCacheIfNeeded() error {
	s.mutex.RLock()
	isExpired := time.Since(s.lastUpdate) > s.cacheTTL
	empty := len(s.assetCache) == 0
	s.mutex.RUnlock()

	if isExpired || empty {
		return s.loadAllAssets()
	}
	return nil
}

// GetByCoinSymbol returns all assets matching the given coin symbol
func (s *AssetService) GetByCoinSymbol(symbol string) ([]staticModels.AssetResponse, error) {
	// Refresh cache if needed
	if err := s.refreshCacheIfNeeded(); err != nil {
		return nil, fmt.Errorf("failed to refresh asset cache: %v", err)
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	symbolUpper := strings.ToUpper(symbol)
	assets, exists := s.assetCache[symbolUpper]
	if !exists {
		return []staticModels.AssetResponse{}, nil // Return empty slice if symbol not found
	}

	return assets, nil
}

// GetAllSymbols returns all available coin symbols
func (s *AssetService) GetAllSymbols() ([]string, error) {
	// Refresh cache if needed
	if err := s.refreshCacheIfNeeded(); err != nil {
		return nil, fmt.Errorf("failed to refresh asset cache: %v", err)
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	symbols := make([]string, 0, len(s.assetCache))
	for symbol := range s.assetCache {
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

// GetAllAssets returns all assets with full details
func (s *AssetService) GetAllAssets() ([]staticModels.AssetResponse, error) {
	// Refresh cache if needed
	if err := s.refreshCacheIfNeeded(); err != nil {
		return nil, fmt.Errorf("failed to refresh asset cache: %v", err)
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var allAssets []staticModels.AssetResponse
	for _, assets := range s.assetCache {
		allAssets = append(allAssets, assets...)
	}

	return allAssets, nil
}

// GetAllAssetsWithPagination returns all assets with pagination and optional blockchain filter
func (s *AssetService) GetAllAssetsWithPagination(limit, offset int, blockchainID string) ([]staticModels.AssetResponse, int, error) {
	// Refresh cache if needed
	if err := s.refreshCacheIfNeeded(); err != nil {
		return nil, 0, fmt.Errorf("failed to refresh asset cache: %v", err)
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var allAssets []staticModels.AssetResponse
	for _, assets := range s.assetCache {
		allAssets = append(allAssets, assets...)
	}

	// Filter by blockchain ID if provided
	if blockchainID != "" {
		var filteredAssets []staticModels.AssetResponse
		for _, asset := range allAssets {
			if asset.BlockchainID == blockchainID {
				filteredAssets = append(filteredAssets, asset)
			}
		}
		allAssets = filteredAssets
	}

	// Sort assets for consistent ordering (by ID) using efficient sort
	sort.Slice(allAssets, func(i, j int) bool {
		// Convert IDs to integers for numeric comparison
		idI, errI := strconv.Atoi(allAssets[i].ID)
		idJ, errJ := strconv.Atoi(allAssets[j].ID)

		// If both are valid numbers, compare numerically
		if errI == nil && errJ == nil {
			return idI < idJ
		}

		// Otherwise fall back to string comparison
		return allAssets[i].ID < allAssets[j].ID
	})

	total := len(allAssets)

	// Apply pagination
	start := offset
	if start > total {
		start = total
	}

	end := start + limit
	if end > total {
		end = total
	}

	paginatedAssets := allAssets[start:end]

	return paginatedAssets, total, nil
}

// ForceRefresh forces a cache refresh
func (s *AssetService) ForceRefresh() error {
	return s.loadAllAssets()
}

// GetCacheStats returns cache statistics
func (s *AssetService) GetCacheStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	totalAssets := 0
	for _, assets := range s.assetCache {
		totalAssets += len(assets)
	}

	return map[string]interface{}{
		"total_symbols":     len(s.assetCache),
		"total_assets":      totalAssets,
		"last_update":       s.lastUpdate,
		"cache_ttl_mins":    int(s.cacheTTL.Minutes()),
		"next_refresh":      s.lastUpdate.Add(s.cacheTTL),
		"total_id_mappings": len(s.idMappings),
		"next_id":           s.nextID,
	}
}

// GenerateAllIDs forces generation of IDs for all assets and saves to file
func (s *AssetService) GenerateAllIDs() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Clear existing cache to force reload
	s.assetCache = make(map[string][]staticModels.AssetResponse)

	// Read all blockchain directories
	entries, err := ioutil.ReadDir(s.assetsPath)
	if err != nil {
		return fmt.Errorf("failed to read assets directory: %v", err)
	}

	totalGenerated := 0
	for _, entry := range entries {
		if entry.IsDir() {
			blockchainName := entry.Name()

			// Skip if blockchain ID mapping doesn't exist
			if s.getBlockchainID(blockchainName) == "" {
				continue
			}

			blockchainPath := filepath.Join(s.assetsPath, blockchainName)

			// Load native token info
			infoPath := filepath.Join(blockchainPath, "info", "info.json")
			if infoData, err := ioutil.ReadFile(infoPath); err == nil {
				var blockchainInfo staticModels.BlockchainInfo
				if err := json.Unmarshal(infoData, &blockchainInfo); err == nil {
					s.getOrCreateAssetID(blockchainName, "", blockchainInfo.Symbol)
					totalGenerated++
				}
			}

			// Load token assets
			assetsDir := filepath.Join(blockchainPath, "assets")
			if entries, err := ioutil.ReadDir(assetsDir); err == nil {
				for _, entry := range entries {
					if entry.IsDir() {
						assetPath := filepath.Join(assetsDir, entry.Name())
						assetInfoPath := filepath.Join(assetPath, "info.json")

						if assetData, err := ioutil.ReadFile(assetInfoPath); err == nil {
							var assetInfo staticModels.AssetInfo
							if err := json.Unmarshal(assetData, &assetInfo); err == nil {
								s.getOrCreateAssetID(blockchainName, entry.Name(), assetInfo.Symbol)
								totalGenerated++
							}
						}
					}
				}
			}
		}
	}

	fmt.Printf("Generated IDs for %d assets. Total mappings: %d\n", totalGenerated, len(s.idMappings))
	return nil
}

// TokenListItem represents a single token from the tokenlist.json
type TokenListItem struct {
	Address string `json:"address"`
	Symbol  string `json:"symbol"`
	Name    string `json:"name"`
}

// TokenList represents the structure of tokenlist.json
type TokenList struct {
	Tokens []TokenListItem `json:"tokens"`
}

// loadSolanaTokenCache loads Solana token symbols from asset files into a static cache
func (s *AssetService) loadSolanaTokenCache() {
	s.tokenCacheMutex.Lock()
	defer s.tokenCacheMutex.Unlock()

	// If already loaded, skip
	if s.tokenCacheLoaded {
		return
	}

	s.solanaTokenCache = make(map[string]string)
	solanaPath := filepath.Join(s.assetsPath, "solana")

	// Load from asset info.json files
	assetsDir := filepath.Join(solanaPath, "assets")
	if entries, err := ioutil.ReadDir(assetsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				mintAddress := entry.Name() // Directory name is the mint address
				assetInfoPath := filepath.Join(assetsDir, mintAddress, "info.json")

				if assetData, err := ioutil.ReadFile(assetInfoPath); err == nil {
					var assetInfo staticModels.AssetInfo
					if err := json.Unmarshal(assetData, &assetInfo); err == nil {
						// Map mint address to symbol (case-insensitive)
						s.solanaTokenCache[strings.ToLower(mintAddress)] = assetInfo.Symbol
					}
				}
			}
		}
	}

	// Also load from tokenlist.json for additional tokens
	tokenListPath := filepath.Join(solanaPath, "tokenlist.json")
	if data, err := ioutil.ReadFile(tokenListPath); err == nil {
		var tokenList TokenList
		if err := json.Unmarshal(data, &tokenList); err == nil {
			// Add tokens from tokenlist (don't override existing ones from info.json)
			for _, token := range tokenList.Tokens {
				normalizedAddr := strings.ToLower(token.Address)
				if _, exists := s.solanaTokenCache[normalizedAddr]; !exists {
					s.solanaTokenCache[normalizedAddr] = token.Symbol
				}
			}
		}
	}

	s.tokenCacheLoaded = true
	fmt.Printf("Loaded %d Solana tokens into cache from asset files\n", len(s.solanaTokenCache))
}

// GetTokenSymbolByMint returns the token symbol for a given Solana mint address
// Uses a static cache loaded from tokenlist.json
func (s *AssetService) GetTokenSymbolByMint(mint string) string {
	// Ensure cache is loaded (happens only once)
	if !s.tokenCacheLoaded {
		s.loadSolanaTokenCache()
	}

	s.tokenCacheMutex.RLock()
	defer s.tokenCacheMutex.RUnlock()

	// Normalize the mint address to lowercase for case-insensitive lookup
	normalizedMint := strings.ToLower(mint)

	// Lookup in cache
	if symbol, exists := s.solanaTokenCache[normalizedMint]; exists {
		return symbol
	}

	// Return the first 4 characters of the mint as fallback
	if len(mint) >= 4 {
		return mint[:4] + "..."
	}
	return "UNKNOWN"
}
