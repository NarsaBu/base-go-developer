package main

import (
	"fmt"
	"time"
)

func main() {
	user1 := NewUser("1", "Tiger", 1000)
	user2 := NewUser("2", "Pigeon", 500)

	go user1.Deposit(200)
	go user1.Deposit(100)
	go user1.Withdraw(40)
	go user1.Withdraw(156)
	go user2.Withdraw(600)
	go user2.Withdraw(100)
	go user2.Withdraw(300)
	go user2.Deposit(50)

	time.Sleep(time.Second)
	fmt.Println(user1)
	fmt.Println(user2)
}
