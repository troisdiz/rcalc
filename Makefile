TARGET=rcalc

fmt:
	@go fmt ./...

lint:
	golangci-lint run ./...


generate: grammar/Rcalc.g4
	antlr -Dlanguage=Go -o rcalc/parser -package parser grammar/Rcalc.g4

compile:
	go build -o bin/$(TARGET) main/main.go

test:
	go test ./rcalc ./main

clean:
	$(RM) -v bin/*
	$(RM) -vrf rcalc/parser
