package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

func guess(totalTries, randomMin, randomMax *int) {
	var history []int
	number := generateRandomNumber(randomMin, randomMax)

	fmt.Printf("Игра 'Угадай число' - от %d до %d началась!\n", *randomMin, *randomMax)
	fmt.Printf("Угадайте число за %d попыток!\n", *totalTries)

	for i := 1; i <= *totalTries; i++ {
		printNumberHistory(&history)

		color.Yellow("Попытка #%d - Введите число: ", i)
		userNumber := scanNumber()
		history = append(history, userNumber)

		if isGuessed(number, userNumber) {
			color.Green("Игра закончена!")
			WriteResultLog("win", i)
			return
		}
	}

	color.Red("Вы проиграли!😢")
	fmt.Println("Секретное число было: ", number)
	WriteResultLog("lose", *totalTries)
}

func generateRandomNumber(randomMin, randomMax *int) int {
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	return randomGenerator.Intn(*randomMax-*randomMin+1) + *randomMin
}

func scanNumber() int {
	var userNumber int

	_, err := fmt.Scan(&userNumber)

	if err != nil {
		fmt.Print("Ошибка при вводе данных. Введите целое число: ")
		scanNumber()
	}

	return userNumber
}

func printNumberHistory(history *[]int) {
	if len(*history) != 0 {
		fmt.Println("Введенные ранее числа: ", *history)
	}
}

func isGuessed(number, userNumber int) bool {
	if number == userNumber {
		fmt.Println("Вы угадали!🙌")
		return true
	}

	if number > userNumber {
		fmt.Println("Секретное число больше👆")
	} else {
		fmt.Println("Секретное число меньше👇")
	}

	printHint(math.Abs(float64(number - userNumber)))
	return false
}

func printHint(value float64) {
	if value <= 5 {
		fmt.Println("🔥 Горячо")
	} else if value <= 15 {
		fmt.Println("🙂 Тепло")
	} else {
		fmt.Println("❄️ Холодно")
	}
}
