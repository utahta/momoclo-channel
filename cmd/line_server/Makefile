.PHONY: all build

all:
	@make build

build:
	GOOS=linux GOARCH=amd64 go build -o line-server

deploy:
	@make build
	scp line-server deploy.line-server:momoclo/line_server/

