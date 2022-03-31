TARGET=rcalc

fmt:
	@go fmt ./...

vet:
	@go vet -v ./...

amd64:
	go build -o bin/$(TARGET) main/main.go

clean:
	$(RM) -v bin/*

