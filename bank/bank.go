package main

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/jeromedesantos/bank/fileops"
)

const accountBalanceFile = "balance.txt"

func main() {
	var accountBalance, err = fileops.GetFloatFromFile(accountBalanceFile)

	if err != nil {
		fmt.Println("ERROR")
		fmt.Println(err)
		fmt.Println("-----------------")
	}

	fmt.Println("Welcome to Go Bank!")
	fmt.Println("Reach us 24/7", randomdata.PhoneNumber())

	for {
		// tanpa import karena masih dalam 1 package (main)
		presentOptions()

		var choice int

		fmt.Print("Your choice: ")
		fmt.Scanln(&choice) // pointer

		switch choice {

		case 1:
			fmt.Println("Your account balance is:", accountBalance)

		case 2:
			fmt.Print("Your deposit: ")
			var depositAmount float64
			fmt.Scanln(&depositAmount)
			if depositAmount <= 0 {
				fmt.Println("Invalid amount. Must be greater than 0.")
				continue
			}
			accountBalance += depositAmount
			fmt.Println("Balance updated! New amount:", accountBalance)
			fileops.WriteFloatToFile(accountBalance, accountBalanceFile)

		case 3:
			var withdrawAmount float64
			fmt.Print("Your withdrawal amount: ")
			fmt.Scanln(&withdrawAmount)

			if withdrawAmount <= 0 {
				fmt.Println("Invalid amount. Must be greater than 0.")
				continue
			}
			if withdrawAmount >= accountBalance {
				fmt.Println("Invalid withdrawal amount!. You can't withdraw more than you have")
				continue
			}

			accountBalance -= withdrawAmount
			fmt.Println("Withdrawal successful! New amount:", accountBalance)
			fileops.WriteFloatToFile(accountBalance, accountBalanceFile)

		default:
			fmt.Println("Goodbye!")
			fmt.Println("Thank you for choosing our bank!")
			return
		}
	}
}
