funcRoot := ./functions

ifeq ($(ENVIRONMENT), production)
	ARCH = GOOS=linux GOARCH=amd64
	LD_FLAGS = -ldflags='-s -w'
	TAGS = -tags=prod
endif

# This is pretty stupid to have to do
ifeq ($(OS), Windows_NT)
	RM_FLAGS = -r -Force
else
	RM_FLAGS = -rf
endif

.PHONY: install dev-run clean compile test

dev-run:
	docker-compose build --no-cache && docker-compose --env-file ./.env up

install:
	go get -u ./... && go mod tidy

clean:
	rm $(RM_FLAGS) $(funcRoot)/bin

compile: clean
	$(ARCH) go build $(LD_FLAGS) $(TAGS) -o $(funcRoot)/bin/server ./cmd/server/main.go

test:
	$(ARCH) go test -v ./...

lint:
	& 'C:\Program Files\Git\bin\sh.exe' --login -c "golangci-lint run -v"