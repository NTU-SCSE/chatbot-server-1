# chatbot-server
Server-side implementation of FYP: SCSE Chatbot.

# Installation
## Golang
Make sure golang is installed, follow the instructions in the following link to install golang:
https://golang.org/doc/install

## chatbot-server
### 1. Clone the repository
`git clone https://github.com/dbakti7/chatbot-server.git`

### 2. Get the following dependencies
```
go get github.com/gorilla/mux
go get github.com/jmoiron/sqlx
go get github.com/marcossegovia/apiai-go
go get github.com/mattn/go-sqlite3
go get github.com/rs/cors
go get github.com/sajari/fuzzy
go get github.com/stretchr/testify/assert
go get github.com/tidwall/gjson
```

### 3. Run the server
```
go run main.go
```
You can check the server is running by accessing:
`localhost:8080`
There should be a "Hello World" message if server is running correctly.
Note: You can change the configuration in `config/config.json`

### 4. Course Parser module
Please go to `courseparser` folder to access the module

### 5. Functional Test
You can test the end-to-end functionality of Dialogflow agent with this functional test.
Please go to `functional_test` folder to access the module.
You can change the target Dialogflow agent at `functional_test/test_config.json`.