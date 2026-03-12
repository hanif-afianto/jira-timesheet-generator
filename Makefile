APP_NAME = jtg
CONFIG_DIR = $(HOME)/.jtg

build:
	go build -o $(APP_NAME) cmd/jtg/main.go
	./$(APP_NAME) install

setup-config:
	mkdir -p $(CONFIG_DIR)
	@if [ ! -f $(CONFIG_DIR)/.env ]; then \
		cp .env.example $(CONFIG_DIR)/.env; \
		echo "Created $(CONFIG_DIR)/.env from .env.example"; \
	else \
		echo "$(CONFIG_DIR)/.env already exists"; \
	fi

run:
	go run cmd/jtg/main.go

clean:
	rm -f $(APP_NAME)

test:
	go test ./...

release:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP_NAME)-darwin-amd64 cmd/jtg/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/$(APP_NAME)-darwin-arm64 cmd/jtg/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux-amd64 cmd/jtg/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/$(APP_NAME)-linux-arm64 cmd/jtg/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME)-windows-amd64.exe cmd/jtg/main.go
	@echo "Binaries created in bin/ directory"
