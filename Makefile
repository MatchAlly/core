SHELL:=/bin/bash

up:
	docker compose up

down:
	docker compose down

seed:
	docker compose up -d db
	go run main.go seed

migrate:
	docker compose up -d db
	go run main.go migrate

init: migrate seed up