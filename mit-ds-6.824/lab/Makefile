GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    BINARY_NAME=mit_ds
    BINARY_UNIX=$(BINARY_NAME)_mac

all: test build

build:
    $(GOBUILD) -o $(BINARY_NAME) -v
test:
    $(GOTEST) -v ./...

test_raft:
    $(GOTEST) -v raft -run $(FOO)

