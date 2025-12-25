# app
export APP_CTX_TIMEOUT=2m
export BCRYPT_COST=12

# миграции
migration_dir=migrations

# сервер
export GIN_MODE=release
export PORT=8080

# nats
nats_port=4222

export NATS_URL=nats://127.0.0.1:$(nats_port)

# postgres
postgres_user=testuser
postgres_password=testpasswd
postgres_port=5432
postgres_db=test
postgres_sslmode=disable

export POSTGRES_DNS=postgresql://$(postgres_user):$(postgres_password)@127.0.0.1:$(postgres_port)/$(postgres_db)?sslmode=$(postgres_sslmode)

# redis
redis_port=6379

export REDIS_ADDR=127.0.0.1:$(redis_port)

run:
	go run cmd/main.go

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=<migration_name>. Example: make migrate-create name=create_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(migration_dir) -seq $(name);

migrate-up:
	migrate -database $(POSTGRES_DNS) -path $(migration_dir) up

migrate-down:
	migrate -database $(POSTGRES_DNS) -path $(migration_dir) down

# Маппинг портов:
#       - 4222 - для клиентов
#       - 6222 - для роутинга
#       - 8222 - для мониторинга
start-nats-server:
	docker run --rm \
		-p 4222:$(nats_port) \
		-p 6222:6222 \
		-p 8222:8222 \
		--name nats-server \
		nats:2.12.3-alpine3.22

start-postgres:
	docker run --rm \
		-e POSTGRES_USER=$(postgres_user) \
		-e POSTGRES_PASSWORD=$(postgres_password) \
		-e POSTGRES_DB=$(postgres_db) \
		-p 5432:$(postgres_port) \
		-v pgdata:/var/lib/postgresql/data \
		--name postgres-db \
		postgres:15.15-alpine3.22

start-redis:
	docker run --rm \
		-p 6379:$(redis_port) \
		--name redis-cache \
		redis:8.4-alpine3.22
