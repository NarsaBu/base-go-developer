package main

import (
	"errors"
	"fmt"
)

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

func (ps *PaymentSystem) ProcessingTransactions(transaction Transaction) error {
	fromUser, fromUseOk := ps.Users[transaction.FromID]
	toUser, toUseOk := ps.Users[transaction.ToID]

	if !fromUseOk || !toUseOk {
		return errors.New("error processing transaction: FromUser or/and ToUser does not exist")
	}

	err := fromUser.Withdraw(transaction.Amount)

	if err != nil {
		return err
	}

	toUser.Deposit(transaction.Amount)
	fmt.Println("Successfully applied transaction: ", transaction)
	return nil
}
