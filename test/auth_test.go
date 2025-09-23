package test

import (
	"github.com/tashunc/nugenesis-wallet-backend/external/data/historical/thrirdParty/moralis"
	"os"
	"testing"
)

func TestHeliusAPI(t *testing.T) {
	os.Setenv("MORALIS_API_KEY", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6Ijc2NTM0MTMwLWU3YTktNGU4NS1hYjIyLWMwYjY5ODhhYjRmMiIsIm9yZ0lkIjoiNDY4NTY0IiwidXNlcklkIjoiNDgyMDMwIiwidHlwZUlkIjoiOGZmMjVhMmUtOTE5OC00NzEwLWIxMzItZTU5OWZhOGFkMmE1IiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE3NTY3MjE0MDIsImV4cCI6NDkxMjQ4MTQwMn0.jYrupIHXLY-2Qvlau5HUdu-qenlBqONZEpy-0SziZT0")
	service := moralis.NewService()
	service.GetAddressInfo("0xe8f04c75b331c309987d161f2ec1bec32170437a")
}
