# Define variables
APP_NAME := forest
SRC_DIR := ./cmd/forest

# Build the application
build:
	go build -o $(APP_NAME) $(SRC_DIR)
	sudo cp $(APP_NAME) /usr/local/bin/$(APP_NAME)

# Run tests
test:
	go test ./...

# Clean up build artifacts
clean:
	rm -f $(APP_NAME)

# Run the application
run: build
	./$(APP_NAME)

.PHONY: build test clean run

