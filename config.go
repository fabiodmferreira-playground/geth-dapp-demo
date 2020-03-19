package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv(filePath string) map[string]string {
	config, err := godotenv.Read(filePath)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return config
}

func loadKeystore(filePath string) *os.File {
	keystore, err := os.Open(filePath)
	if err != nil {
		log.Fatal(
			"could not load keystore from location %s: %v\n",
			filePath,
			err,
		)
	}

	return keystore
}

// updateEnvFile updates our env file with a key-value pair
func updateEnvFile(config map[string]string, k string, val string) {

	config[k] = val

	err := godotenv.Write(config, ENV_FILE)
	if err != nil {
		log.Printf("failed to update %s: %v\n", ENV_FILE, err)
	}
}
