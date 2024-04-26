lint:
	golangci-lint run -v --timeout=2m

up:
	docker-compose -f ./infra/docker-compose.yml up -d

down:
	docker-compose -f ./infra/docker-compose.yml down -v

doc:
	# http://localhost:8080/pkg/github.com/soltanat/metrics/?m=all
	godoc -http=:8080

