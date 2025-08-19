BIN ?= taskmansm7
PKG ?= .
GOBIN ?= $(HOME)/.local/bin
BUILD_DIR ?= dist
LDFLAGS ?= -s -w

.DEFAULT_GOAL := install

.PHONY: build install uninstall reinstall run fmt vet test tidy clean build-linux build-all completion-bash completion-zsh completion-fish

build:
	mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS)" -trimpath -o $(BUILD_DIR)/$(BIN) $(PKG)

install:
	mkdir -p $(GOBIN)
	go build -ldflags="$(LDFLAGS)" -trimpath -o $(GOBIN)/$(BIN) $(PKG)

uninstall:
	rm -f $(GOBIN)/$(BIN)

reinstall:
	$(MAKE) uninstall
	$(MAKE) install

run:
	go run .

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

tidy:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR)

build-linux:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -trimpath -o $(BUILD_DIR)/$(BIN)-linux-amd64 $(PKG)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -trimpath -o $(BUILD_DIR)/$(BIN)-linux-arm64 $(PKG)

completion-bash:
	$(GOBIN)/$(BIN) completion bash | sudo tee /etc/bash_completion.d/$(BIN) >/dev/null

completion-zsh:
	mkdir -p ~/.zsh/completions
	$(GOBIN)/$(BIN) completion zsh > ~/.zsh/completions/_$(BIN)

completion-fish:
	mkdir -p ~/.config/fish/completions
	$(GOBIN)/$(BIN) completion fish > ~/.config/fish/completions/$(BIN).fish
