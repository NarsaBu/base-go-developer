package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex

type PaymentSystem struct {
	Users        map[string]*User
	Transactions []Transaction
}

func NewPaymentSystem() *PaymentSystem {
	return &PaymentSystem{
		Users:        make(map[string]*User),
		Transactions: make([]Transaction, 0),
	}
}

func (ps *PaymentSystem) AddUser(user User) {
	mu.Lock()
	defer mu.Unlock()

	ps.Users[user.ID] = &user
	fmt.Printf("Added user %s to the payment system\n", user.ID)
}

func (ps *PaymentSystem) AddTransaction(transaction Transaction) {
	mu.Lock()
	defer mu.Unlock()

	ps.Transactions = append(ps.Transactions, transaction)
	fmt.Printf("Added transaction to the payment system\n")
}

func (ps *PaymentSystem) ProcessingTransactions() {
	mu.Lock()
	defer mu.Unlock()

	for _, transaction := range ps.Transactions {
		fromUser, fromUseOk := ps.Users[transaction.FromID]
		toUser, toUseOk := ps.Users[transaction.ToID]

		if !fromUseOk || !toUseOk {
			println("Error processing transaction: FromUser or/and ToUser does not exist")
			continue
		}

		if err := fromUser.Withdraw(transaction.Amount); err == nil {
			toUser.Deposit(transaction.Amount)
			fmt.Println("Successfully applied transaction: ", transaction)
		} else {
			println("Error processing transaction: ", err)
		}
	}

	ps.Transactions = ps.Transactions[:0]
}
