TARGET=rcalc

fmt:
	@go fmt ./...

vet:
	@go vet -v ./...

compile:
	go build -o bin/$(TARGET) main/main.go

test:
	go test ./rcalc ./main

clean:
	$(RM) -v bin/*

