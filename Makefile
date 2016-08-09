.PHONY: install test

install:
	@glide install

serve:
	@goapp serve app

test:
	@go test -v ./crawler/...
