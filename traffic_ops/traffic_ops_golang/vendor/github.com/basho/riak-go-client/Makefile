.PHONY: all install-deps
.PHONY: unit-test integration-test security-test test fmt help
.PHONY: fmt lint protogen release

PROJDIR = $(realpath $(CURDIR))
PROTOC_VERSION := $(shell protoc --version)

export RIAK_HOST = localhost
export RIAK_PORT = 8087

all: install-deps lint test

install-deps:
	cd $(PROJDIR) && go get -t github.com/basho/riak-go-client/...

lint: install-deps
	cd $(PROJDIR) && go tool vet -shadow=true -shadowstrict=true $(PROJDIR)
	cd $(PROJDIR) && go vet github.com/basho/riak-go-client/...

unit-test: lint
	cd $(PROJDIR) && go test -v

integration-test: lint
	cd $(PROJDIR) && go test -v -tags='integration timeseries'

timeseries-test: lint
	cd $(PROJDIR) && go test -v -tags=timeseries

security-test: lint
	cd $(PROJDIR) && go test -v -tags=security

test: integration-test

fmt:
	cd $(PROJDIR) && gofmt -s -w .

protogen:
ifeq ($(PROTOC_VERSION),)
	$(error The protoc command is required to parse proto files)
endif
ifneq ($(PROTOC_VERSION),libprotoc 2.6.1)
	$(error protoc must be version 2.6.1)
endif
	$(PROJDIR)/build/protogen $(PROJDIR)

release:
ifeq ($(VERSION),)
	$(error VERSION must be set to deploy this code)
endif
ifeq ($(RELEASE_GPG_KEYNAME),)
	$(error RELEASE_GPG_KEYNAME must be set to deploy this code)
endif
	@$(PROJDIR)/tools/build/publish $(VERSION) master validate
	@git tag --sign -a "$(VERSION)" -m "riak-go-client $(VERSION)" --local-user "$(RELEASE_GPG_KEYNAME)"
	@git push --tags
	@$(PROJDIR)/tools/build/publish $(VERSION) master 'Riak Go Client' 'riak-go-client'

help:
	@echo ''
	@echo ' Targets:'
	@echo '----------------------------------------------------------'
	@echo ' all                  - Run everything                    '
	@echo ' fmt                  - Format code                       '
	@echo ' lint                 - Run "go vet"                      '
	@echo ' test                 - Run unit & integration tests      '
	@echo ' unit-test            - Run unit tests                    '
	@echo ' integration-test     - Run integration tests             '
	@echo ' timeseries-test      - Run timeseries tests              '
	@echo ' security-test        - Run security tests                '
	@echo '----------------------------------------------------------'
	@echo ''
