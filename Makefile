.PHONY: install fmt test serve-app deploy-app-dev deploy-app-prod

install:
	@glide install

fmt:
	goimports -w ./appengine

test:
	@goapp test -v -race ./appengine/model/...
	@goapp test -v -race ./appengine/lib/crawler/...

serve-app:
	@make -C appengine/app serve

deploy-app-prod:
	@make -C appengine/app deploy-prod

deploy-app-dev:
	@make -C appengine/app deploy-dev

