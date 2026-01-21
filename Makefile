APP_NAME = jtg

build:
	go build -o $(APP_NAME) cmd/jtg/main.go
	./$(APP_NAME) install

run:
	go run cmd/jtg/main.go

clean:
	rm -f $(APP_NAME)

test:
	go test ./...
