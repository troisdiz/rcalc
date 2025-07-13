TARGET=bin/rcalc

FIND=find
# Idea from https://stackoverflow.com/a/12099167/32015
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	FIND=gfind
endif

SRCS=$(shell $(FIND) rcalc main cmd -iname '*.go' -type f)

GRAMMAR_SRC=grammar/Rcalc.g4
GRAMMAR_OUTPUT_DIR=rcalc/parser
GRAMMAR_WITNESS=$(GRAMMAR_OUTPUT_DIR)/Rcalc.tokens

PROTO_BASE_NAME=rcalc-stack
PROTO_SRC=proto/$(PROTO_BASE_NAME).proto
PROTO_GO_FILE=rcalc/protostack/$(PROTO_BASE_NAME).pb.go

all: $(TARGET)

fmt:
	@go fmt ./...

lint: $(GRAMMAR_WITNESS)
	golangci-lint run ./...

$(GRAMMAR_WITNESS): $(GRAMMAR_SRC)
	antlr -Dlanguage=Go -Werror -o $(GRAMMAR_OUTPUT_DIR) -Xexact-output-dir -package parser grammar/Rcalc.g4

$(PROTO_GO_FILE): $(PROTO_SRC)
	mkdir -p rcalc/protostack
	protoc -I=. --go_opt=module=troisdizaines.com/rcalc --go_out=rcalc $<

$(TARGET): $(GRAMMAR_WITNESS) $(SRCS) $(PROTO_GO_FILE)
	go build -o $(TARGET) main/main.go

generate: $(GRAMMAR_WITNESS) $(PROTO_GO_FILE)

compile: $(TARGET)

test: $(GRAMMAR_WITNESS) $(PROTO_GO_FILE)
	go test ./rcalc ./main

run: $(TARGET)
	$(TARGET)

clean:
	$(RM) -vrf bin
	$(RM) -vrf $(GRAMMAR_OUTPUT_DIR)
	$(RM) -vrf rcalc/protostack

.PHONY: all test compile generate clean
