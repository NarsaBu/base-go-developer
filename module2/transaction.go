package main

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

func NewTransaction(fromId, toId string, amount float64) *Transaction {
	return &Transaction{
		FromID: fromId,
		ToID:   toId,
		Amount: amount,
	}
}
