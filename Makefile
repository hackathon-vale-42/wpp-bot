CMD_DIR = cmd

BIN_DIR = bin

DOCKER_COMPOSE = docker-compose.yaml


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

.PHONY: docker
docker:
	@docker compose -f $(DOCKER_COMPOSE) up -d

.PHONY: nuke
nuke:
	@-docker rm -f $$(docker ps -aq)
	@-docker network rm $$(docker network ls -q)
	@-docker volume rm $$(docker volume ls -q)
