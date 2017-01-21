.PHONY: install fmt test serve-app deploy-app-dev deploy-app-prod

install:
	@glide install

fmt:
	gofmt -w ./appengine ./linenotify ./log ./twitter
	goimports -w ./appengine ./linenotify ./log ./twitter

test:
	@go test -v -race ./appengine/model/...
	@go test -v -race ./appengine/lib/crawler/...
	@go test -v -race ./linenotify/...

serve-app:
	@cd appengine/app && make serve

deploy-app-prod:
	@cd appengine/app && make deploy-prod

deploy-app-dev:
	@cd appengine/app && make deploy-dev

