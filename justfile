dev:
	npx nodemon -x "~/go/bin/templ generate && go run ./cmd/server/main.go" -e .go,.html --signal SIGINT
