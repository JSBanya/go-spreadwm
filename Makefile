BINARY=spreadwm
export GOPATH=$(CURDIR)/vendor/:$(CURDIR)/cmd/

all: build

build:
	go fmt $(CURDIR)/cmd/*.go
	go build -o $(BINARY) $(CURDIR)/cmd/*.go

fetch:
	go get golang.org/x/image/bmp

clean:
	@go clean
	@$(RM) -rf vendor
	@$(RM) $(BINARY)

.PHONY: all build fetch clean
