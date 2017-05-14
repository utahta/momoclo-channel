.PHONY: install fmt test serve-app deploy-app-dev deploy-app-prod

install:
	@glide install

fmt:
	goimports -w ./appengine ./linenotify ./log ./twitter

test:
	@goapp test -v -race ./appengine/model/...
	@goapp test -v -race ./appengine/lib/crawler/...
	@goapp test -v -race ./linenotify/...

serve-app:
	@cd appengine/app && make serve

deploy-app-prod:
	@cd appengine/app && make deploy-prod

deploy-app-dev:
	@cd appengine/app && make deploy-dev

