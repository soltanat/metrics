lint:
	golangci-lint run -v --timeout=2m

up:
	docker-compose -f ./infra/docker-compose.yml up -d

down:
	docker-compose -f ./infra/docker-compose.yml down -v
