up:
	docker-compose up -d --build

down:
	docker-compose down

test:
	go test -v -cover ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	del coverage.out