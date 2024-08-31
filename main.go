package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	d "gitto/db"

	a "gitto/assitant"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	assistant := a.LLMAssistant{}
	database := d.InitDB().Assign()
	for {
		fmt.Print("> ")
		scanner.Scan()
		userInput := strings.ToLower(scanner.Text())

		if userInput == "exit" || userInput == "quit" || userInput == "bye" {
			fmt.Println("Goodbye!")
			break
		}
		assistant.Process(database, userInput)

	}
}
