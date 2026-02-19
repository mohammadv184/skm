GOCMD:=go
GOBUILD:=$(GOCMD) build
GOTEST:=$(GOCMD) test

EXE =
ifeq (windows,$(shell go env GOOS))
	EXE = .exe
endif

.PHONY: build
# build builds all binaries
build: $(wildcard ./cmd/*)
	$(info ******************** building binaries ********************)
	@for dir in $^; do \
		$(GOBUILD) -o build/$$(basename $$dir) -v $(ARGS) ./$$dir; \
	done
	$(info binaries are built in ./build directory)

.PHONY: clean
# clean cleans all build and vendor directories
clean:
	$(info ******************** cleaning binaries ********************)
	@rm -rf ./build
	@rm -rf ./vendor
	@rm -rf ./assets/mandocs
	@rm -rf ./assets/auto-completions


.PHONY: test
# test runs all tests
test:
	$(info ******************** running tests ********************)
	$(GOTEST) -timeout 30m -p 1 $$(go list ./... | grep -v "vendor/")


.PHONY: lint
# lint runs all linters
lint:
ifeq (, $(shell which golangci-lint))
	$(error "could not find golangci-lint in $(PATH), see: https://golangci-lint.run/docs/welcome/install/local for installation instructions")
else
	$(info ******************** running lint tools ********************)
	@golangci-lint run ./...
endif


.PHONY: lint-fix
# lint-fix runs all linters and fixes all fixable issues
lint-fix:
ifeq (, $(shell which golangci-lint))
	$(error "could not find golangci-lint in $(PATH), see: https://golangci-lint.run/docs/welcome/install/local for installation instructions")
else
	$(info ******************** running lint tools and fixing issues ********************)
	@golangci-lint run ./... --fix
endif

.PHONY: vendor
# vendor vendors all dependencies
vendor:
	$(info ******************** vendoring dependencies ********************)
	go mod vendor

.PHONY: manpages
# manpages generates manpages for all commands
manpages:
	$(info ******************** generating manpages ********************)
	@rm -rf ./assets/mandocs
	@mkdir -p ./assets/mandocs
	@go run ./cmd/skm mandoc | gzip -c -9 >./assets/mandocs/skm.1.gz
	$(info manpages are generated in ./man directory)

.PHONY: auto-completions
# auto-completions generates auto-completions for all supported shells
auto-completions:
	$(info ******************** generating auto-completions ********************)
	@rm -rf ./assets/auto-completions
	@mkdir -p ./assets/auto-completions
	@for sh in bash zsh fish powershell; do \
    	go run ./cmd/skm completion "$$sh" >"./assets/auto-completions/skm.$$sh"; \
    done