run:
	go run main.go -local

up:
	docker-compose up -d --build

.PHONY: run up