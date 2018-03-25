# chatbot-server
Server-side implementation of FYP: SCSE Chatbot.

# Installation
## Golang
Make sure golang is installed, follow the instructions in the following link to install golang:
https://golang.org/doc/install

## chatbot-server
1. Clone the repository
`git clone https://github.com/dbakti7/chatbot-server.git`

2. Get the following dependencies
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

3. Run the server
`go run main.go`
You can check the server is running by accessing:
`localhost:8080`
There should be a "Hellow World" message if server is running correctly.