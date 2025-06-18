BINARY        ?= appsync
GITHUB_TOKEN  ?= $(shell gh auth token)
TENANTS_ROOT  ?= /path/to/dpe/tenants/.applications
SAMPLE_ROOT   ?= ./sample
REPOS_FILE    ?= repos.yaml

.PHONY: all build install sync-ott generate sync-sample test-features clean

all: build

build:
	go build -o $(BINARY) .

install: build
	chmod +x $(BINARY)

sync-ott: install
	./$(BINARY) sync \
		--root       $(TENANTS_ROOT) \
		--repos-file $(REPOS_FILE) \
		--mode       feature \
		--token      "$(GITHUB_TOKEN)"

generate: install
	./$(BINARY) generate \
		--root       $(SAMPLE_ROOT) \
		--repos-file $(REPOS_FILE) \
		--token      "$(GITHUB_TOKEN)"

sync-sample: install
	./$(BINARY) sync \
		--root       $(SAMPLE_ROOT) \
		--repos-file $(REPOS_FILE) \
		--mode       feature \
		--token      "$(GITHUB_TOKEN)"

test-features:
	@go test ./features -v

clean:
	rm -f $(BINARY)
