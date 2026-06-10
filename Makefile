include .env
export

export PROJECT_ROOT=$(shell pwd)

.DEFAULT_GOAL := help

env-up: ## env: Запустить окружение проекта (Postgres, Redis)
	@docker compose up -d postgres redis

env-down: ## env: Остановить окружение проекта
	@docker compose down

env-cleanup: ## env: Очистить данные БД и Redis. Опасность утери данных!
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down && \
		rm -rf ${PROJECT_ROOT}/out/pgdata ${PROJECT_ROOT}/out/redisdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forward: ## env: Открыть порты сервисов окружения (Postgres: 5432, Redis: 6379)
	@docker compose up -d postgres-port-forwarder redis-port-forwarder

env-port-close: ## env: Закрыть порты сервисов окружения
	@docker compose down postgres-port-forwarder redis-port-forwarder

logs-cleanup: ## env: Очистить файлы логов из out/logs
	@read -p "Очистить все log файлы? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

swagger-gen: ## env: Сгенерировать актуальную Swagger спецификацию через Docker
	@docker compose run --rm swagger \
		init \
		-g cmd/app/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

ps: ## env: Посмотреть запущенные Docker Compose сервисы
	@docker compose ps

migrate-create: ## PostgreSQL: Создать новую версию схемы данных (Пример: make migrate-create seq=init_schema)
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init_schema"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up: ## PostgreSQL: Накатить миграции
	@make migrate-action action=up

migrate-down: ## PostgreSQL: Откатить миграции
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable \
		"$(action)"

app-run: ## Golang приложение: Запустить локально на хост-системе (для локальной разработки)
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export DB_HOST=localhost && \
	export REDIS_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/app/main.go

app-deploy: ## Golang приложение: Запустить в Docker Compose (для деплоя)
	@docker compose up -d --build app

app-undeploy: ## Golang приложение: Остановить Docker Compose сервис приложения
	@docker compose down app

test: ## Тестирование: Запустить все юнит и интеграционные тесты
	@go test -v ./...

help: ## Показать справку по командам
	@echo "=== Центр управления проектом Музыкального Сервиса ==="
	@echo ""
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
