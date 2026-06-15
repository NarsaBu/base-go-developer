package main

import (
	"errors"
	"sync"
)

type User struct {
	ID      string
	Name    string
	Balance float64
	mu      sync.Mutex
}

func NewUser(id, name string, balance float64) *User {
	return &User{
		ID:      id,
		Name:    name,
		Balance: balance,
	}
}

func (u *User) Deposit(value float64) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.Balance += value
}

func (u *User) Withdraw(value float64) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if value > u.Balance {
		println("Can`t perform withdraw: withdraw value is larger then balance")
		return errors.New("balance is lower then withdraw value")
	}

	u.Balance -= value
	return nil
}
