.PHONY: install test build-prots serve deploy-dev deploy-prod

install:
	@glide install

test:
	@go test -v ./crawler/...

build-protos:
	@protoc line/protos/line.proto --go_out=plugins=grpc:.

serve-app:
	@cd appengine/app && make serve

deploy-app-prod:
	@cd appengine/app && make deploy-prod

deploy-app-dev:
	@cd appengine/app && make deploy-dev
