package test

import (
	helius_models "github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/helius"
	"os"
	"testing"
)

func TestHeliusAPI(t *testing.T) {
	os.Setenv("HELIUS_API_KEY", "8f33306f-f9aa-4f07-aac3-8489d1ade5ab")
	service := helius_models.NewService()
	service.GetAddressInfo("DDws22Z91d3ZzxPFCqvh1BWZY1zyZzLzGHVXXQw5bhwc", 10, 0)
}
