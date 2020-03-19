package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fabiodmferreira/geth-dapp-demo/quiz"
)

const ErrTransactionWait = "if you've just started the application, wait a while for the network to confirm your transaction."

func DialEthereum(gateway string) *ethclient.Client {
	client, err := ethclient.Dial(gateway)
	if err != nil {
		log.Fatalf("could not connect to Ethereum gateway: %v\n", err)
	}

	return client
}

func NewSession(ctx context.Context, keystore *os.File, keystorepass string) *quiz.QuizSession {
	auth, err := bind.NewTransactor(keystore, keystorepass)
	if err != nil {
		log.Fatal("%s\n", err)
	}

	// Return session without contract instance
	return &quiz.QuizSession{
		TransactOpts: *auth,
		CallOpts: bind.CallOpts{
			From:    auth.From,
			Context: ctx,
		},
	}
}

// NewContract deploys a contract if no existing contract exists
func NewContract(session *quiz.QuizSession, client *ethclient.Client, question string, answer string) string {
	// Hash answer before sending it over Ethereum network.
	contractAddress, tx, instance, err := quiz.DeployQuiz(&session.TransactOpts, client, question, stringToKeccak256(answer))
	if err != nil {
		log.Fatalf("could not deploy contract: %v\n", err)
	}
	fmt.Printf("Contract deployed! Wait for tx %s to be confirmed.\n", tx.Hash().Hex())

	session.Contract = instance
	return contractAddress.Hex()
}

// LoadContract loads a contract if one exists
func LoadContract(session *quiz.QuizSession, client *ethclient.Client, contractAddr string) *quiz.QuizSession {

	addr := common.HexToAddress(contractAddr)
	instance, err := quiz.NewQuiz(addr, client)
	if err != nil {
		log.Fatalf("could not load contract: %v\n", err)
		log.Println(ErrTransactionWait)
	}
	session.Contract = instance
	return session
}

// readQuestion prints out question stored in contract.
func readQuestion(session *quiz.QuizSession) {
	qn, err := session.Question()
	if err != nil {
		log.Printf("could not read question from contract: %v\n", err)
		log.Println(ErrTransactionWait)
		return
	}
	fmt.Printf("Question: %s\n", qn)
	return
}

// sendAnswer sends answer to contract as a keccak256 hash.
func sendAnswer(session *quiz.QuizSession, ans string) {
	// Send answer
	txSendAnswer, err := session.SendAnswer(stringToKeccak256(ans))
	if err != nil {
		log.Printf("could not send answer to contract: %v\n", err)
		return
	}
	fmt.Printf("Answer sent! Please wait for tx %s to be confirmed.\n", txSendAnswer.Hash().Hex())
	return
}

// checkCorrect makes a contract message call to check if
// the current account owner has answered the question correctly.
func checkCorrect(session *quiz.QuizSession) {
	win, err := session.CheckBoard()
	if err != nil {
		log.Printf("could not check leaderboard: %v\n", err)
		log.Println(ErrTransactionWait)
		return
	}
	fmt.Printf("Were you correct?: %v\n", win)
	return
}

// Utility functions

// stringToKeccak256 converts a string to a keccak256 hash of type [32]byte
func stringToKeccak256(s string) [32]byte {
	var output [32]byte
	copy(output[:], crypto.Keccak256([]byte(s))[:])
	return output
}
