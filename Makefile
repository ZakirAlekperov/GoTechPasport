.PHONY: help build run test test-coverage lint clean

# Переменные
APP_NAME=techpassport
BIN_DIR=bin
CMD_DIR=cmd/techpassport
GO=go
GOFLAGS=-v

# Цвета для вывода
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## Показать помощь
	@echo "${GREEN}Доступные команды:${NC}"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${YELLOW}%-20s${NC} %s\n", $$1, $$2}'

build: ## Сборка приложения
	@echo "${GREEN}Сборка ${APP_NAME}...${NC}"
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)
	@echo "${GREEN}✓ Сборка завершена: $(BIN_DIR)/$(APP_NAME)${NC}"

build-release: ## Сборка release версии
	@echo "${GREEN}Сборка release версии...${NC}"
	@mkdir -p $(BIN_DIR)
	$(GO) build -ldflags="-s -w" -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)
	@echo "${GREEN}✓ Release сборка завершена${NC}"

run: ## Запуск приложения
	@echo "${GREEN}Запуск ${APP_NAME}...${NC}"
	$(GO) run ./$(CMD_DIR)

test: ## Запуск тестов
	@echo "${GREEN}Запуск тестов...${NC}"
	$(GO) test -v ./...

test-coverage: ## Запуск тестов с покрытием
	@echo "${GREEN}Запуск тестов с покрытием...${NC}"
	$(GO) test -v -coverprofile=coverage.txt -covermode=atomic ./...
	$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "${GREEN}✓ Отчет сохранен: coverage.html${NC}"

test-race: ## Запуск тестов с проверкой race conditions
	@echo "${GREEN}Запуск тестов с -race...${NC}"
	$(GO) test -race -v ./...

lint: ## Запуск линтера
	@echo "${GREEN}Запуск линтера...${NC}"
	@which golangci-lint > /dev/null || (echo "${RED}golangci-lint не установлен. Установите: https://golangci-lint.run/usage/install/${NC}" && exit 1)
	golangci-lint run ./...

fmt: ## Форматирование кода
	@echo "${GREEN}Форматирование кода...${NC}"
	$(GO) fmt ./...
	@echo "${GREEN}✓ Форматирование завершено${NC}"

vet: ## Проверка кода go vet
	@echo "${GREEN}Проверка go vet...${NC}"
	$(GO) vet ./...

mod-download: ## Загрузка зависимостей
	@echo "${GREEN}Загрузка зависимостей...${NC}"
	$(GO) mod download

mod-tidy: ## Очистка зависимостей
	@echo "${GREEN}Очистка зависимостей...${NC}"
	$(GO) mod tidy

clean: ## Очистка артефактов сборки
	@echo "${GREEN}Очистка...${NC}"
	rm -rf $(BIN_DIR)
	rm -f coverage.txt coverage.html
	@echo "${GREEN}✓ Очистка завершена${NC}"

install-tools: ## Установка инструментов разработки
	@echo "${GREEN}Установка инструментов...${NC}"
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "${GREEN}✓ Инструменты установлены${NC}"

.DEFAULT_GOAL := help
