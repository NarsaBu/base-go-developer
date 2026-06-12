package main

import (
	"fmt"
	"sync"
)

func main() {
	paymentSystem := initData()
	ch := make(chan Transaction, len(paymentSystem.Transactions))
	var wg sync.WaitGroup

	addTransactionsToChannel(paymentSystem, ch)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go Worker(ch, &wg, paymentSystem)
	}

	close(ch)
	wg.Wait()

	fmt.Println(paymentSystem.Users["1"])
	fmt.Println(paymentSystem.Users["2"])
}

func addTransactionsToChannel(ps *PaymentSystem, ch chan Transaction) {
	for _, transaction := range ps.Transactions {
		ch <- transaction
	}
}

func Worker(ch chan Transaction, wg *sync.WaitGroup, ps *PaymentSystem) {
	defer wg.Done()

	for transaction := range ch {
		if err := ps.ProcessingTransactions(transaction); err != nil {
			println(err)
		}
	}
}

func initData() *PaymentSystem {
	paymentSystem := NewPaymentSystem()

	user1 := NewUser("1", "Tiger", 1000)
	user2 := NewUser("2", "Pigeon", 500)

	paymentSystem.AddUser(*user1)
	paymentSystem.AddUser(*user2)

	transaction1 := NewTransaction("1", "2", 200)
	transaction2 := NewTransaction("2", "1", 50)

	paymentSystem.AddTransaction(*transaction1)
	paymentSystem.AddTransaction(*transaction2)

	return paymentSystem
}
