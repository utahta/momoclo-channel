.PHONY: install test build-protos serve-app deploy-app-dev deploy-app-prod

install:
	@glide install

fmt:
	gofmt -w .
	goimports -w .

test:
	@go test -v -race ./crawler/... ./ustream/... ./appengine/lib/util/... ./appengine/model/...

build-protos:
	@protoc linebot/protos/linebot.proto --go_out=plugins=grpc:.

serve-app:
	@cd appengine/app && make serve

deploy-app-prod:
	@cd appengine/app && make deploy-prod

deploy-app-dev:
	@cd appengine/app && make deploy-dev

deploy-linebot:
	@cd cmd/linebot_server && make deploy
