dev:
	~/go/bin/templ generate --watch &
	npx nodemon@3.0.1 -x "go run ./cmd/server/main.go" -e .go,.html --signal SIGINT
