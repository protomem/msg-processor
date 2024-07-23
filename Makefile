PROJECT := msg-processor

.DEFAULT_GOAL := run


.PHONY: build
build:
	go build -o ./bin/${PROJECT} .


.PHONY: run
run: cfg_file=.env.dev
run: build
	./bin/${PROJECT} -cfg=${cfg_file}


.PHONY: test
test:
	go test -v ./...

