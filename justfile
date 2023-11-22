dev:
	npx nodemon -x "go run ./cmd/server/main.go" -e .go,.html --signal SIGINT
