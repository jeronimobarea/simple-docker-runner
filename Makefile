run:
	docker-compose up --build

test:
	docker-compose up --build -d
	go test ./...
