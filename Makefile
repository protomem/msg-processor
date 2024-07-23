PROJECT := msg-processor

.DEFAULT_GOAL := run


.PHONY: build
build:
	go build -o ./bin/${PROJECT} .


.PHONY: run
run: build
	./tmp/${PROJECT}


.PHONY: test
test:
	go test -v ./...

