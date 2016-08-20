.PHONY: install test

install:
	@glide install

build-prots:
	@protoc grpc/line/protos/line.proto --go_out=plugins=grpc:.

serve:
	@ln -nfs .env.local app/.env
	@goapp serve app

test:
	@go test -v ./crawler/...
