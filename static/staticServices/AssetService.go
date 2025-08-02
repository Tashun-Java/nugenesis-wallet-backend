package staticServices

import (
	"encoding/json"
	"fmt"
	"github.com/tashunc/nugenesis-wallet-backend/static/staticModels"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// AssetService handles blockchain assets management
type AssetService struct {
	assetCache map[string][]staticModels.AssetResponse
	mutex      sync.RWMutex
	lastUpdate time.Time
	cacheTTL   time.Duration
	assetsPath string
}

// NewAssetService creates a new AssetService instance
func NewAssetService() *AssetService {
	return &AssetService{
		assetCache: make(map[string][]staticModels.AssetResponse),
		cacheTTL:   30 * time.Minute, // Cache for 30 minutes
		assetsPath: "assets/blockchains",
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
				Symbol:     blockchainInfo.Symbol,
				Name:       blockchainInfo.Name,
				Blockchain: blockchainName,
				Address:    "", // Native tokens don't have contract addresses
				Type:       blockchainInfo.Type,
				Decimals:   blockchainInfo.Decimals,
				Status:     blockchainInfo.Status,
			}

			if logoExists {
				nativeAsset.LogoPath = fmt.Sprintf("assets/blockchains/%s/info/logo.png", blockchainName)
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
							Symbol:     assetInfo.Symbol,
							Name:       assetInfo.Name,
							Blockchain: blockchainName,
							Address:    entry.Name(), // Directory name is the contract address/asset ID
							Type:       assetInfo.Type,
							Decimals:   assetInfo.Decimals,
							Status:     assetInfo.Status,
						}

						if logoExists {
							tokenAsset.LogoPath = fmt.Sprintf("assets/blockchains/%s/assets/%s/logo.png", blockchainName, entry.Name())
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
		"total_symbols":  len(s.assetCache),
		"total_assets":   totalAssets,
		"last_update":    s.lastUpdate,
		"cache_ttl_mins": int(s.cacheTTL.Minutes()),
		"next_refresh":   s.lastUpdate.Add(s.cacheTTL),
	}
}
