# Geth DApp Demo

Golang/Ethereum API experiment

## Features

- Load environment variables from file;
- Connect to ethereum blockchain;
- Create Smart Contract;
- Execute Smart Contract methods.

## Getting Started

1. Set up Rinkeby testnet endpoint on Infura.io
2. Create new account `geth --datadir . account new`
3. Update `.env`file
```
ANSWER="<quiz_answer>"
GATEWAY="<ethereum_endpoint_url>"
KEYSTORE="<keystore_file_path>"
KEYSTOREPASS="<keystore_passphrase>"
QUESTION="<quiz_question>"
```
4. Compile smart contract
```
$ npm install solc@0.5.15
$ node_modules/.bin/solc --abi --bin ./quiz/quiz.sol -o build
```
5. Generate go binding file from abi and bin files
```
$ abigen --abi="build/__quiz_quiz_sol_Quiz.abi" --bin="build/__quiz_quiz_sol_Quiz.bin" --pkg=quiz --out="quiz/quiz.go"
```
6. Run program `go run main`

## Files and functions

- config.go
  - loadEnv
  - loadKeystore
  - updateEnvFile
- ethereum.go
  - DialEthereum
  - NewSession
  - NewContract
  - LoadContract
  - readQuestion
  - sendAnswer
  - checkCorrect
  - stringToKeccak256
- main.go
  - main
  - readStringStdin

