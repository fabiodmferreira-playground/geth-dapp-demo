package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	GATEWAY      = "GATEWAY"
	CONTRACTADDR = "CONTRACTADDR"
	ANSWER       = "ANSWER"
	QUESTION     = "QUESTION"
	KEYSTORE     = "KEYSTORE"
	KEYSTOREPASS = "KEYSTOREPASS"
)

const ENV_FILE = ".env"

func main() {
	config := loadEnv(ENV_FILE)

	ctx := context.Background()

	client := DialEthereum(config[GATEWAY])
	defer client.Close()

	// accountAddress := common.HexToAddress("0x608Fe39257edab277d31487dBD276c43c8610211")
	// balance, _ := client.BalanceAt(ctx, accountAddress, nil)
	// fmt.Printf("Balance: %d\n", balance)

	keystore := loadKeystore(config[KEYSTORE])
	defer keystore.Close()

	session := NewSession(ctx, keystore, config[KEYSTOREPASS])

	contractAddr := config[CONTRACTADDR]

	if contractAddr == "" {
		contractAddress := NewContract(session, client, config[QUESTION], config[ANSWER])
		updateEnvFile(config, CONTRACTADDR, contractAddress)
	} else {
		LoadContract(session, client, contractAddr)
	}

	for {
		fmt.Printf(
			"Pick an option:\n" + "" +
				"1. Show question.\n" +
				"2. Send answer.\n" +
				"3. Check if you answered correctly.\n" +
				"4. Exit.\n",
		)

		// Reads a single UTF-8 character (rune)
		// from STDIN and switches to case.
		switch readStringStdin() {
		case "1":
			readQuestion(session)
			break
		case "2":
			fmt.Println("Type in your answer")
			sendAnswer(session, readStringStdin())
			break
		case "3":
			checkCorrect(session)
			break
		case "4":
			fmt.Println("Bye!")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
			break
		}
	}
}

// Utility functions

// readStringStdin reads a string from STDIN and strips any trailing \n characters from it.
func readStringStdin() string {
	reader := bufio.NewReader(os.Stdin)
	inputVal, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("invalid option: %v\n", err)
		return ""
	}

	output := strings.TrimSuffix(inputVal, "\n") // Important!
	return output
}
