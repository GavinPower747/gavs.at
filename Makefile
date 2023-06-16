ifeq ($(CI), true)
	CI_ARCH = GOOS=linux GOARCH=amd64
endif

compile:
	$(CI_ARCH) go build -o ./functions/bin/server ./cmd/server/main.go

test:
	$(CI_ARCH) go test -v ./...

dev-run: 
	docker-compose --env-file ./.env up