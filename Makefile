# golang1.9 or latest
# 1. make help
# 2. make dep
# 3. make build
# ...

SRC := github.com/harrylee2015/monitor/
APP := build/monitor
LDFLAGS := -ldflags "-w -s"
PKG_LIST := `go list ./... | grep -v "vendor" | grep -v "test" | grep -v "mocks"
BUILD_FLAGS = -ldflags
.PHONY: default build release linter race test fmt protobuf clean help

default: build

build: ## Build the binary file
	@go build $(BUILD_FLAGS) -v -i -o  $(APP) $(SRC)
	@cp monitor.toml build/
	@cp install.sql build/
	@cp initDataBase.sh build/

release: ## Build the binary file
	@go build -v -i -o $(APP) $(LDFLAGS) $(SRC) 
	@cp monitor.toml build/
	@cp install.sql build/
	@cp initDataBase.sh build/

linter: ## Use gometalinter check code, ignore some unserious warning
	@res=$$(gometalinter.v2 -t --sort=linter --enable-gc --deadline=2m --disable-all \
	--enable=gofmt \
	--enable=gosimple \
	--enable=deadcode \
	--enable=unconvert \
	--enable=interfacer \
	--enable=varcheck \
	--enable=structcheck \
	--enable=goimports \
	--vendor ./...) \

	if [ -n "$$res" ]; then \
		echo "$${res}"; \
		exit 1; \
		fi;
	@find . -name '*.sh' -not -path "./vendor/*" | xargs shellcheck

race: ## Run data race detector
	@go test -race -short $(PKG_LIST)

test: ## Run unittests
	@go test -race $(PKG_LIST)

fmt: fmt_proto fmt_shell ## go fmt
	@go fmt ./...
	@find . -name '*.go' -not -path "./vendor/*" | xargs goimports -l -w

.PHONY: fmt_proto fmt_shell
fmt_proto: ## go fmt protobuf file
	@find . -name '*.proto' -not -path "./vendor/*" | xargs clang-format -i

clean: ## Remove previous build
	@rm -rf $(shell find . -name 'datadir' -not -path "./vendor/*")
	@rm -rf build/monitor*
	@rm -rf build/*.log
	@rm -rf build/logs
	@go clean

protobuf: ## Generate protbuf file of types package
	@protoc --go_out=plugins=grpc:./  types/*.proto

cleandata:
	rm -rf build/datadir

.PHONY: checkgofmt
checkgofmt: ## get all go files and run go fmt on them
	@files=$$(find . -name '*.go' -not -path "./vendor/*" | xargs gofmt -l -s); if [ -n "$$files" ]; then \
		  echo "Error: 'make fmt' needs to be run on:"; \
		  echo "${files}"; \
		  exit 1; \
		  fi;
	@files=$$(find . -name '*.go' -not -path "./vendor/*" | xargs goimports -l -w); if [ -n "$$files" ]; then \
		  echo "Error: 'make fmt' needs to be run on:"; \
		  echo "${files}"; \
		  exit 1; \
		  fi;
