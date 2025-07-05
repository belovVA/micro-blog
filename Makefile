.PHONY: build-up test cover

build-up:
	#docker compose up -d
	bash start.sh
test:
#	go clean -testcache
#	go test -covermode=atomic -coverpkg=$(COVERPKG) -coverprofile=coverage.out $(PKGS)

cover:
#	go tool cover -func=coverage.out