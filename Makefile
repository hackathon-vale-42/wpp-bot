CMD_DIR = cmd

BIN_DIR = bin


.PHONY: all
all: build

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

.PHONY: build
build: $(BIN_DIR)
	@go build -o $(BIN_DIR) ./$(CMD_DIR)/...

.PHONY: run
run: build
	./$(BIN_DIR)/wpp-bot
