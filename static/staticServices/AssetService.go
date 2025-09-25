package staticServices

import (
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/external/models/general"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticModels"
	"io/ioutil"
	"os"
	"path/filepath"
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
	assetCache    map[string][]staticModels.AssetResponse
	idMappings    map[string]string           // asset_key -> id
	nextID        int                         // counter for generating new IDs
	blockchainMap map[string]general.CoinType // blockchain name -> CoinType ID
	mutex         sync.RWMutex
	lastUpdate    time.Time
	cacheTTL      time.Duration
	assetsPath    string
	idMappingFile string
}

// NewAssetService creates a new AssetService instance
func NewAssetService() *AssetService {
	service := &AssetService{
		assetCache:    make(map[string][]staticModels.AssetResponse),
		idMappings:    make(map[string]string),
		blockchainMap: make(map[string]general.CoinType),
		nextID:        1,
		cacheTTL:      30 * time.Minute, // Cache for 30 minutes
		assetsPath:    "./assets/blockchains",
		idMappingFile: "./assets/id_mappings.json",
	}
	service.initializeBlockchainMapping()
	service.loadIDMappings()
	return service
}

// initializeBlockchainMapping creates the mapping between blockchain names and CoinType IDs
func (s *AssetService) initializeBlockchainMapping() {
	s.blockchainMap["ethereum"] = general.Ethereum
	s.blockchainMap["polygon"] = general.Polygon
	s.blockchainMap["aurora"] = general.Aurora
	s.blockchainMap["avalanchec"] = general.AvalancheCChain
	s.blockchainMap["binance"] = general.Binance
	s.blockchainMap["base"] = general.Base
	s.blockchainMap["blast"] = general.Blast
	s.blockchainMap["cfxevm"] = general.ConfluxeSpace
	s.blockchainMap["kcc"] = general.KuCoinCommunityChain
	s.blockchainMap["linea"] = general.Linea
	s.blockchainMap["manta"] = general.MantaPacific
	s.blockchainMap["meter"] = general.Meter
	s.blockchainMap["oasis"] = general.Oasis
	s.blockchainMap["optimismgoerli"] = general.OptimismTest
	s.blockchainMap["scroll"] = general.Scroll
	s.blockchainMap["sepolia"] = general.Sepolia
	s.blockchainMap["smartchain"] = general.SmartChain
	s.blockchainMap["solana"] = general.Solana
	s.blockchainMap["terra"] = general.Terra
	s.blockchainMap["zklink"] = general.ZkLinkNova
	s.blockchainMap["zksync"] = general.Zksync
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

			nativeAsset := staticModels.AssetResponse{
				ID:         s.getOrCreateAssetID(blockchainName, "", blockchainInfo.Symbol),
				Symbol:     blockchainInfo.Symbol,
				Name:       blockchainInfo.Name,
				Blockchain: blockchainName,
				Address:    "", // Native tokens don't have contract addresses
				Type:       blockchainInfo.Type,
				Decimals:   blockchainInfo.Decimals,
				Status:     blockchainInfo.Status,
			}

			// Add blockchain ID if mapping exists
			if blockchainID := s.getBlockchainID(blockchainName); blockchainID != "" {
				nativeAsset.BlockchainID = blockchainID
			}

			if logoExists {
				nativeAsset.LogoPath = fmt.Sprintf("/api/static/assetsLogo/blockchains/%s/info/logo.png", blockchainName)
			}

			assets = append(assets, nativeAsset)
		}
	}

	// Load token assets
	assetsDir := filepath.Join(blockchainPath, "assets")
	if entries, err := ioutil.ReadDir(assetsDir); err == nil {
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
							ID:         s.getOrCreateAssetID(blockchainName, entry.Name(), assetInfo.Symbol),
							Symbol:     assetInfo.Symbol,
							Name:       assetInfo.Name,
							Blockchain: blockchainName,
							Address:    entry.Name(), // Directory name is the contract address/asset ID
							Type:       assetInfo.Type,
							Decimals:   assetInfo.Decimals,
							Status:     assetInfo.Status,
						}

						// Add blockchain ID if mapping exists
						if blockchainID := s.getBlockchainID(blockchainName); blockchainID != "" {
							tokenAsset.BlockchainID = blockchainID
						}

						if logoExists {
							tokenAsset.LogoPath = fmt.Sprintf("/api/static/assetsLogo/blockchains/%s/assets/%s/logo.png", blockchainName, entry.Name())
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
