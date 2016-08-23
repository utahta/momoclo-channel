.PHONY: install test build-prots serve deploy-dev deploy-prod rollback

install:
	@glide install

build-protos:
	@protoc grpc/line/protos/line.proto --go_out=plugins=grpc:.

serve:
	@cp app/.env.local app/env
	@goapp serve app

test:
	@go test -v ./crawler/...

deploy-prod:
	@cp app/.env.prod app/env
	@appcfg.py -A momoclo-channel update app

deploy-dev:
	@cp app/.env.dev app/env
	@appcfg.py -A momoclo-channel-dev update app
