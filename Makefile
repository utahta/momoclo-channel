.PHONY: install test build-prots serve deploy

install:
	@glide install

build-prots:
	@protoc grpc/line/protos/line.proto --go_out=plugins=grpc:.

serve:
	@cp app/.env.local app/env
	@goapp serve app

test:
	@go test -v ./crawler/...

deploy:
	@cp app/.env.prod app/env
	@appcfg.py -A momoclo-channel update app
