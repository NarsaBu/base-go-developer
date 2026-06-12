package main

import (
	"sync"
)

var mu sync.Mutex

type User struct {
	ID      string
	Name    string
	Balance float64
}

func NewUser(id, name string, balance float64) *User {
	return &User{
		ID:      id,
		Name:    name,
		Balance: balance,
	}
}

func (u *User) Deposit(value float64) {
	mu.Lock()
	defer mu.Unlock()

	u.Balance += value
}

func (u *User) Withdraw(value float64) {
	mu.Lock()
	defer mu.Unlock()

	if value > u.Balance {
		println("Can`t perform withdraw: withdraw value is larger then balance")
		return
	}

	u.Balance -= value
}
