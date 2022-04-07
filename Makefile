TARGET=rcalc

fmt:
	@go fmt ./...

lint:
	golangci-lint run ./...

compile:
	go build -o bin/$(TARGET) main/main.go

test:
	go test ./rcalc ./main

clean:
	$(RM) -v bin/*

