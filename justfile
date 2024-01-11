watch-templ:
	npx nodemon@3.0.1 -x "~/go/bin/templ generate" -e .templ

dev:
	# just watch-templ & ~/go/bin/wgo run cmd/server/main.go && fg
	templ generate --watch & wgo run cmd/server/main.go
