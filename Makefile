lint:
	@echo ? Running golangci-lint
	@docker-compose -f tools/golang/docker-compose.yml up lint

unit:
	@echo ? Running golangci-lint
	@docker-compose -f tools/golang/docker-compose.yml up unit
