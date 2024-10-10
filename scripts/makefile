# Variables for the Go compiler and the output binary names
PRODUCER_BINARY = ../cmd/producer_app/producer_app.exe
CONSUMER_BINARY = ../cmd/consumer_app/consumer_app.exe

# Default target: Build both programs
all: build

# Build both producer and consumer binaries
build: $(PRODUCER_BINARY) $(CONSUMER_BINARY)

# Compile producer service
$(PRODUCER_BINARY): ../cmd/producer_app/main.go
	go build -o $(PRODUCER_BINARY) ../cmd/producer_app

# Compile consumer service
$(CONSUMER_BINARY): ../cmd/consumer_app/main.go
	go build -o $(CONSUMER_BINARY) ../cmd/consumer_app

# Run both services in parallel-something goes wrong !!! do not use (modified path in createType)
run:
	@echo "Starting producer and consumer services..."
	go run ../cmd/producer_app/main.go
	go run ../cmd/consumer_app/main.go

