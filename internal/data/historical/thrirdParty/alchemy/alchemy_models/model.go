package alchemy_models

type TokensByAddressRequest struct {
	Addresses []AddressRequest `json:"addresses"`
}

type AddressRequest struct {
	Address  string   `json:"address"`
	Networks []string `json:"networks"`
}

type TokensByAddressResponse struct {
	Data TokenData `json:"data"`
}

type TokenData struct {
	Tokens  []TokenInfo `json:"tokens"`
	PageKey *string     `json:"pageKey"`
}

type TokenInfo struct {
	Address       string        `json:"address"`
	Network       string        `json:"network"`
	TokenAddress  *string       `json:"tokenAddress"`
	TokenBalance  string        `json:"tokenBalance"`
	TokenMetadata TokenMetadata `json:"tokenMetadata"`
	TokenPrices   []TokenPrice  `json:"tokenPrices"`
}

type TokenMetadata struct {
	Symbol   *string `json:"symbol"`
	Decimals *int    `json:"decimals"`
	Name     *string `json:"name"`
	Logo     *string `json:"logo"`
}

type TokenPrice struct {
	Currency      string `json:"currency"`
	Value         string `json:"value"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
}
