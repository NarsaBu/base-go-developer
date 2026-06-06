package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var totalTries, randomMin, randomMax int

func init() {
	color.NoColor = false
	chooseDifficulty()
}

func main() {
	for {
		guess(&totalTries, &randomMin, &randomMax)
		retryGame()
	}

}

func chooseDifficulty() {
	var difficulty string
	text := `Выберите сложность:
		- Easy: 1–50, 15 попыток
		- Medium: 1–100, 10 попыток
		- Hard: 1–200, 5 попыток`

	fmt.Println(text)

	fmt.Scan(&difficulty)

	switch strings.ToLower(difficulty) {
	case "easy":
		totalTries, randomMin, randomMax = 15, 1, 50
	case "medium":
		totalTries, randomMin, randomMax = 10, 1, 100
	case "hard":
		totalTries, randomMin, randomMax = 5, 1, 200
	default:
		fmt.Println("Неизвестный формат. Введите Easy/Medium/Hard")
		chooseDifficulty()
	}
}

func retryGame() {
	var wantToRetry string

	fmt.Println("Сыграть еще раз?. Y/n")
	fmt.Scan(&wantToRetry)

	switch strings.ToLower(wantToRetry) {
	case "y":
		chooseDifficulty()
	case "n":
		os.Exit(0)
	default:
		fmt.Println("Неизвестный формат. Введите Y/n")
		retryGame()
	}
}
