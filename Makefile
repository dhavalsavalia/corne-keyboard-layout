.PHONY: build run clean

# Build the flash TUI
build:
	go build -o flash ./cmd/flash

# Run the flash TUI
run: build
	./flash

# Clean build artifacts
clean:
	rm -f flash
	go clean
