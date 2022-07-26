TARGET=bin/rcalc

FIND=find
# Idea from https://stackoverflow.com/a/12099167/32015
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	FIND=gfind
endif

SRCS=$(shell $(FIND) rcalc main -iname '*.go' -type f)

GRAMMAR_SRC=grammar/Rcalc.g4
GRAMMAR_OUTPUT_DIR=rcalc/parser
GRAMMAR_WITNESS=$(GRAMMAR_OUTPUT_DIR)/grammar/Rcalc.tokens

all: $(TARGET)

fmt:
	@go fmt ./...

lint: $(GRAMMAR_WITNESS)
	golangci-lint run ./...

$(GRAMMAR_WITNESS): $(GRAMMAR_SRC)
	antlr -Dlanguage=Go -o $(GRAMMAR_OUTPUT_DIR) -package parser grammar/Rcalc.g4

$(TARGET): $(GRAMMAR_WITNESS) $(SRCS)
	go build -o $(TARGET) main/main.go

generate: $(GRAMMAR_WITNESS)

compile: $(TARGET)

test: $(GRAMMAR_WITNESS)
	go test ./rcalc ./main

clean:
	$(RM) -vrf bin
	$(RM) -vrf $(GRAMMAR_OUTPUT_DIR)

.PHONY: all test compile generate clean
