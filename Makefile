.PHONY: build clean

GO=CGO_ENABLED=0 GO111MODULE=on go
GOFLAGS=-ldflags  "-s -w"

SERVICES=driver-example
build:
	$(GO) build $(GOFLAGS) -o $(SERVICES) main.go

clean:
	rm -f $(SERVICES)

