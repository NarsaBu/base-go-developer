package main

import (
	"fmt"
	"time"
)

func main() {
	paymentSystem := NewPaymentSystem()

	user1 := NewUser("1", "Tiger", 1000)
	user2 := NewUser("2", "Pigeon", 500)

	paymentSystem.AddUser(*user1)
	paymentSystem.AddUser(*user2)

	transaction1 := NewTransaction("1", "2", 200)
	transaction2 := NewTransaction("2", "1", 50)

	paymentSystem.AddTransaction(*transaction1)
	paymentSystem.AddTransaction(*transaction2)

	go paymentSystem.ProcessingTransactions()

	time.Sleep(time.Second)
	fmt.Println(paymentSystem.Users["1"])
	fmt.Println(paymentSystem.Users["2"])
}
