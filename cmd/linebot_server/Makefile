.PHONY: all build

all:
	@make build

build:
	GOOS=linux GOARCH=amd64 go build -o linebot-server

deploy:
	@make build
	scp linebot-server deploy.linebot-server:momoclo/linebot_server/
