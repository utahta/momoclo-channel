.PHONY: install fmt test serve-app deploy-app-dev deploy-app-prod

install:
	@glide install

fmt:
	gofmt -w .
	goimports -w .

test:
	@go test -v -race ./appengine/model/...

serve-app:
	@cd appengine/app && make serve

deploy-app-prod:
	@cd appengine/app && make deploy-prod

deploy-app-dev:
	@cd appengine/app && make deploy-dev

