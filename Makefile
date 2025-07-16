.PHONY: build-up test cover

build-up:
	#docker compose up -d
	bash start.sh
test:
	go clean -testcache
	go test -race -bench=. ./internal/service/...

cover:
#	go tool cover -func=coverage.out