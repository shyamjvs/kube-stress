export GOPROXY=direct
UNAME_S = $(shell uname -s)
GO_INSTALL_FLAGS=-ldflags="-s -w"
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
TARGET := kube-stress

build: fmt
	go mod tidy
ifeq ($(UNAME_S), Darwin)
	time GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(TARGET) $(GO_INSTALL_FLAGS) $V
else
	time GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(TARGET) $(GO_INSTALL_FLAGS) $V
endif

fmt:
	@gofmt -l -w $(SRC)

clean:
	@rm kube-stress

.PHONY: build fmt clean
