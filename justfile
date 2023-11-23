dev:
	~/go/bin/templ generate --watch &
	npx nodemon -x "go run ./cmd/server/main.go" -e .go,.html --signal SIGINT
