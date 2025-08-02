package staticServices

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type NonceRecord struct {
	Hash      string    `json:"hash"`
	Used      bool      `json:"used"`
	Timestamp time.Time `json:"timestamp"`
}

type NonceStore struct {
	nonces map[string]*NonceRecord
	mutex  sync.RWMutex
}

func NewNonceStore() *NonceStore {
	return &NonceStore{
		nonces: make(map[string]*NonceRecord),
	}
}

// Register a nonce hash for future use
func (ns *NonceStore) RegisterNonce(nonce string, timestamp time.Time) string {
	// Create hash of nonce + timestamp for validation
	hasher := sha256.New()
	hasher.Write([]byte(nonce + strconv.FormatInt(timestamp.Unix(), 10)))
	hash := hex.EncodeToString(hasher.Sum(nil))

	ns.mutex.Lock()
	ns.nonces[hash] = &NonceRecord{
		Hash:      hash,
		Used:      false,
		Timestamp: timestamp,
	}
	ns.mutex.Unlock()

	return hash
}

// Validate and consume a nonce
func (ns *NonceStore) ValidateAndConsumeNonce(nonce string, timestamp time.Time, providedHash string) bool {
	// Recreate expected hash
	hasher := sha256.New()
	hasher.Write([]byte(nonce + strconv.FormatInt(timestamp.Unix(), 10)))
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	// Check if provided hash matches expected
	if expectedHash != providedHash {
		return false
	}

	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	record, exists := ns.nonces[expectedHash]
	if !exists {
		return false
	}

	// Check if already used
	if record.Used {
		return false
	}

	// Check timestamp validity (within 5 minutes)
	if time.Since(timestamp) > 5*time.Minute {
		delete(ns.nonces, expectedHash)
		return false
	}

	// Mark as used
	record.Used = true

	// Remove after use (prevents replay)
	delete(ns.nonces, expectedHash)

	return true
}

func (ns *NonceStore) CleanupExpired() {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	cutoff := time.Now().Add(-5 * time.Minute)
	for hash, record := range ns.nonces {
		if record.Timestamp.Before(cutoff) {
			delete(ns.nonces, hash)
		}
	}
}

// Middleware for client-generated nonce authentication
func ClientNonceAuthMiddleware(store *NonceStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		nonce := c.GetHeader("X-Nonce")
		timestampStr := c.GetHeader("X-Nonce-Timestamp")
		hash := c.GetHeader("X-Nonce-Hash")

		if nonce == "" || timestampStr == "" || hash == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing nonce headers (X-Nonce, X-Nonce-Timestamp, X-Nonce-Hash required)",
			})
			c.Abort()
			return
		}

		// Parse timestamp
		timestampInt, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid timestamp format"})
			c.Abort()
			return
		}

		timestamp := time.Unix(timestampInt, 0)

		if !store.ValidateAndConsumeNonce(nonce, timestamp, hash) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid, expired, or already used nonce",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
