package main

import (
	"errors"
	"fmt"
	"os"
)

func main () {
	revenue, err1 := getUserInput("Revenue: ")
	expenses, err2 := getUserInput("Expenses: ")
	taxRate, err3 := getUserInput("Tax Rate: ")

	if err1 != nil || err2 != nil || err3 != nil {
		panic(err1)
	}

	ebt, profit, ratio := calculateFinancials(revenue, expenses, taxRate)

	writeToFile(ebt, profit, ratio)
}

func writeToFile(ebt, profit, ratio float64) error {
	file, err := os.Create("financials.txt")
	if err != nil {
		return errors.New("Error creating file")
	}

	results := fmt.Sprintf("EBT: %.1f\nProfit: %.1f\nRatio: %.3f\n", ebt, profit, ratio)
	file.WriteString(results)

	fmt.Println("--------- Profit Calculation ---------")
	fmt.Print(results)
	return nil
}

// func storeResults(ebt, profit, ratio float64) {
// 	results := fmt.Sprintf("EBT: %.1f\nProfit: %.1f\nRatio: %.3f\n", ebt, profit, ratio )

// 	os.WriteFile("results.txt",  []byte(results), 0644)
// }

func calculateFinancials(revenue, expenses, taxRate float64) (ebt float64, profit float64, ratio float64) {
	ebt = revenue - expenses
	profit = ebt * (1 - taxRate/100)
	ratio = ebt / profit
	return
}

func getUserInput(infoText string) (float64, error) {
	var userInput float64

	fmt.Print(infoText)
	fmt.Scan(&userInput)

	if userInput <= 0 {
		return 0, errors.New("Value must be a positive number")
	}

	return userInput, nil
}