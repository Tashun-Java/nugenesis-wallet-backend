package models

type Transaction struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	Token   string `json:"token"`
	Amount  string `json:"amount"`
	Value   string `json:"value"`
	Address string `json:"address"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Fee     string `json:"fee"`
	Hash    string `json:"hash"`
}
