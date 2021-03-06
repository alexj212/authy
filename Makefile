#!make
-include .env
export $(shell [ -f ".env" ] && sed 's/=.*//' .env)

export BIN_DIR=./bin
export RELEASE_DIR=./release
export PROJ_PATH=github.com/alexj212/authy
export APP_NAME=authy


export DATE := $(shell date +%Y.%m.%d-%H%M)
export BUILT_ON_IP := $(shell [ $$(uname) = Linux ] && hostname -i || hostname )
export RUNTIME_VER := $(shell go version)
export BUILT_ON_OS=$(shell uname -a)

export LATEST_COMMIT := $(shell git log --pretty=format:'%h' -n 1 2> /dev/null)
export COMMIT_CNT := $(shell git rev-list --all 2> /dev/null | wc -l | sed 's/ //g' )
export BRANCH := $(shell git branch  2> /dev/null |grep -v "no branch"| grep \*|cut -d ' ' -f2)

ifeq ($(BRANCH),)
BRANCH := master
endif

ifeq ($(COMMIT_CNT),)
COMMIT_CNT := 0
endif

export BUILD_NUMBER := ${BRANCH}-${COMMIT_CNT}

export COMPILE_LDFLAGS=-s -X "main.BuildDate=${DATE}" \
                          -X "main.LatestCommit=${LATEST_COMMIT}" \
                          -X "main.BuildNumber=${BUILD_NUMBER}" \
                          -X "main.BuiltOnIP=${BUILT_ON_IP}" \
                          -X "main.BuiltOnOs=${BUILT_ON_OS}" \
						  -X "main.RuntimeVer=${RUNTIME_VER}"





build_info: check_prereq ## Build the container
	@echo ''
	@echo '---------------------------------------------------------'
	@echo 'DATE              $(DATE)'
	@echo 'LATEST_COMMIT     $(LATEST_COMMIT)'
	@echo 'BUILD_NUMBER      $(BUILD_NUMBER)'
	@echo 'BUILT_ON_IP       $(BUILT_ON_IP)'
	@echo 'BUILT_ON_OS       $(BUILT_ON_OS)'
	@echo 'RUNTIME_VER       $(RUNTIME_VER)'

	@echo 'BRANCH            $(BRANCH)'
	@echo 'COMMIT_CNT        $(COMMIT_CNT)'
	@echo 'COMPILE_LDFLAGS   $(COMPILE_LDFLAGS)'

	@echo 'PATH              $(PATH)'
	@echo '---------------------------------------------------------'
	@echo ''



####################################################################################################################
##
## help for each task - https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
##
####################################################################################################################
.PHONY: help

help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help



####################################################################################################################
##
## Build of binaries
##
####################################################################################################################
all: fmt check authy ## build authy and run tests

binaries: authy ## build binaries in bin dir

create_dir:
	@mkdir -p $(BIN_DIR)

check_prereq: create_dir

build_app: create_dir
	go build -o $(BIN_DIR)/$(BIN_NAME) -a -ldflags '$(COMPILE_LDFLAGS)' $(APP_PATH)

	## GOOS=linux GOARCH=arm   go build -o bin/main-linux-arm main.go
	## GOOS=linux GOARCH=arm64 go build -o bin/main-linux-arm64 main.go


authy: ## build_info ## build authy binary in bin dir
	@echo "build authy"
	make BIN_NAME=authy APP_PATH=$(PROJ_PATH) build_app
	@echo ''
	@echo ''


build_release:
	@mkdir -p $(RELEASE_DIR)
	@echo ""
	@echo ""
	@echo "building release: $(RELEASE_DIR)/$(BIN_NAME)-$(GOOS)-$(GOARCH)$(BIN_EXTENSION)"
	@echo ""
	BIN_EXTENSION=$(BIN_EXTENSION) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(RELEASE_DIR)/$(BIN_NAME)-$(GOOS)-$(GOARCH)$(BIN_EXTENSION) -a -ldflags '$(COMPILE_LDFLAGS)' $(APP_PATH)


release_artifacts: ## build_info fmt check  ## build release binaries into release dir
	@rm -rf $(RELEASE_DIR)
	@mkdir -p $(RELEASE_DIR)
	@echo "build release artifacts"
	make BIN_EXTENSION=     BIN_NAME=authy GOOS=linux   GOARCH=amd64 APP_PATH=$(PROJ_PATH) build_release
	@echo ''
	make BIN_EXTENSION=.exe BIN_NAME=authy GOOS=windows GOARCH=amd64 APP_PATH=$(PROJ_PATH) build_release
	@echo ''
	make BIN_EXTENSION=     BIN_NAME=authy GOOS=darwin  GOARCH=amd64 APP_PATH=$(PROJ_PATH) build_release
	@echo ''
	@echo ''

####################################################################################################################
##
## Cleanup of binaries
##
####################################################################################################################

clean: clean_authy  ## clean all binaries in bin dir


clean_binary: ## clean binary in bin dir
	rm -f $(BIN_DIR)/$(BIN_NAME)

clean_authy: ## clean dumper
	make BIN_NAME=authy clean_binary



test: ## run tests
	go test -v $(PROJ_PATH)

fmt: ## run fmt on project
	#go fmt $(PROJ_PATH)/...
	gofmt -s -d -w -l .

doc: ## launch godoc on port 6060
	godoc -http=:6060

deps: ## display deps for project
	go list -f '{{ join .Deps  "\n"}}' . |grep "/" | grep -v $(PROJ_PATH)| grep "\." | sort |uniq

lint: ## run lint on the project
	golint ./...

staticcheck: ## run staticcheck on the project
## staticcheck -ignore "$(shell cat .checkignore)" .
	staticcheck .

vet: ## run go vet on the project
	go vet .

tools: ## install dependent tools for code analysis
	go get -u honnef.co/go/tools
	go get -u github.com/gordonklaus/ineffassign
	go get -u github.com/fzipp/gocyclo
	go get -u golang.org/x/lint/golint

gocyclo: ## run gocyclo on the project
	@ gocyclo -over 15 $(shell find . -name "*.go" |egrep -v "pb\.go|_test\.go")

check: staticcheck gocyclo ## run code checks on the project
