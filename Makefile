start:
	docker-compose up -d

start-build:
	docker-compose up -d --build

stop:
	docker-compose down

restart:
	make stop && make start

restart-build:
	make stop && make start-build

# dev:
# 	air
dev-watch:
	APP_SERVICE=api air -c .air.toml

dev:
	APP_SERVICE=api go run main.go

cronjob:
	APP_SERVICE=cronjob go run main.go

install:
	export GOPRIVATE=gitlab.finema.co/finema/* && git config --global url."git@gitlab.finema.co:".insteadOf "https://gitlab.finema.co/" && go get

logs:
	 docker logs -f api

seed:
	APP_SERVICE=seed go run main.go

# migrate:
# 	docker-compose up -d --build migration

migrate:
	APP_SERVICE=migration go run main.go

test:
	go test ./...
