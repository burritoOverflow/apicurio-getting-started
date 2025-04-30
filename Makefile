.PHONY: all build up down clean

DOCKER_COMPOSE = docker-compose
GO = go
BIN_DIR = ./bin

CREATE_ARTIFACT_SRC = ./client/create-artifact.go
CREATE_ARTIFACT_BIN= $(BIN_DIR)/create-artifacts
SEARCH_ARTIFACT_SRC = ./client/search-artifacts.go
SEARCH_ARTIFACT_BIN = $(BIN_DIR)/search-artifacts

all: build up

build: $(CREATE_ARTIFACT_BIN) $(SEARCH_ARTIFACT_BIN)

$(CREATE_ARTIFACT_BIN): $(CREATE_ARTIFACT_SRC)
	@echo "Building create-artifact..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $@ $(CREATE_ARTIFACT_SRC)

$(SEARCH_ARTIFACT_BIN): $(SEARCH_ARTIFACT_SRC)
	@echo "Building search-artifacts..."
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $@ $(SEARCH_ARTIFACT_SRC)

up:
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d

down:
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

clean:
	@echo "Cleaning up binaries..."
	@rm -rf $(BIN_DIR)
