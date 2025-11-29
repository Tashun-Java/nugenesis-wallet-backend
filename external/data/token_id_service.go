package data

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

// TokenIDMapping represents a single mapping entry from the JSON file
type TokenIDMapping struct {
	AssetKey string `json:"asset_key"`
	ID       string `json:"id"`
}

// TokenIDService provides token ID lookup functionality
type TokenIDService struct {
	mappings map[string]string
	mu       sync.RWMutex
	loaded   bool
}

var (
	tokenIDService *TokenIDService
	once           sync.Once
)

// GetTokenIDService returns the singleton instance of TokenIDService
func GetTokenIDService() *TokenIDService {
	once.Do(func() {
		tokenIDService = &TokenIDService{
			mappings: make(map[string]string),
		}
	})
	return tokenIDService
}

// LoadMappings loads the ID mappings from the JSON file
func (s *TokenIDService) LoadMappings(filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If already loaded, skip
	if s.loaded {
		return nil
	}

	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read ID mappings file: %w", err)
	}

	// Parse the JSON
	var mappings []TokenIDMapping
	if err := json.Unmarshal(data, &mappings); err != nil {
		return fmt.Errorf("failed to parse ID mappings JSON: %w", err)
	}

	// Build the lookup map
	for _, mapping := range mappings {
		s.mappings[strings.ToLower(mapping.AssetKey)] = mapping.ID
	}

	s.loaded = true
	return nil
}

// GetTokenID retrieves the token ID for a given chain and token address
// Returns empty string if not found
func (s *TokenIDService) GetTokenID(chain, tokenAddress string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Normalize the chain name to match the mapping format
	normalizedChain := normalizeChainName(chain)

	// Build the asset key in the format: chain-tokenAddress
	assetKey := fmt.Sprintf("%s-%s", normalizedChain, tokenAddress)

	// Lookup with case-insensitive key
	if id, exists := s.mappings[strings.ToLower(assetKey)]; exists {
		return id
	}

	return ""
}

// GetTokenIDForNative retrieves the token ID for a native token using chain and symbol
// Returns empty string if not found
func (s *TokenIDService) GetTokenIDForNative(chain, symbol string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Normalize the chain name to match the mapping format
	normalizedChain := normalizeChainName(chain)

	// Build the asset key in the format: chain-SYMBOL-native
	assetKey := fmt.Sprintf("%s-%s-native", normalizedChain, symbol)

	// Lookup with case-insensitive key
	if id, exists := s.mappings[strings.ToLower(assetKey)]; exists {
		return id
	}

	return ""
}

// normalizeChainName converts chain names to the format used in id_mappings.json
func normalizeChainName(chain string) string {
	// Map common chain names to their asset key format
	chainMapping := map[string]string{
		"eth":    "ethereum",
		"matic":  "polygon",
		"pol":    "polygon",
		"sol":    "solana",
		"bsc":    "smartchain",
		"bnb":    "smartchain",
		"avax":   "avalanchec",
		"avaxc":  "avalanchec",
		"ftm":    "fantom",
		"arb":    "arbitrum",
		"op":     "optimism",
		"base":   "base",
		"zksync": "zksync",
		"linea":  "linea",
		"scroll": "scroll",
		"blast":  "blast",
		"mantle": "mantle",
		"sei":    "sei",
		"sonic":  "sonic",
		"metis":  "metis",
		"celo":   "celo",
		"ronin":  "ronin",
		"opbnb":  "opbnb",
		"zeta":   "zetachain",
		"tron":   "tron",
		"trx":    "tron",
	}

	normalized := strings.ToLower(chain)
	if mapped, exists := chainMapping[normalized]; exists {
		return mapped
	}

	return normalized
}
