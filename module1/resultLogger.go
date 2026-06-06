package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ResultTable struct {
	Results []Result `json:"results"`
}

type Result struct {
	Date       time.Time `json:"date"`
	Exodus     string    `json:"exodus"`
	TriesCount int       `json:"triesCount"`
}

const filename = "results.json"

func NewResultTable() ResultTable {
	return ResultTable{Results: []Result{}}
}

func writeDataToFile(rt *ResultTable) {
	data, _ := json.MarshalIndent(rt, "", "  ")

	os.WriteFile(filename, data, 0644)
}

func loadResultTable() *ResultTable {
	data, err := os.ReadFile(filename)

	if os.IsNotExist(err) {
		fmt.Println("Файл не найден. Выполняется инициализация")
		writeDataToFile(new(NewResultTable()))
	}

	var resultTable ResultTable

	json.Unmarshal(data, &resultTable)

	return &resultTable
}

func WriteResultLog(exodus string, triesCount int) {
	resultTable := loadResultTable()

	result := Result{
		Date:       time.Now(),
		Exodus:     exodus,
		TriesCount: triesCount,
	}

	resultTable.Results = append(resultTable.Results, result)

	writeDataToFile(resultTable)
}
