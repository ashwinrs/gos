run:
	go generate ./...
	go run cmd/petstore.go

start-db:
	brew services start postgresql

stop-db:
	brew services stop postgresql