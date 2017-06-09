install:
	@dep ensure

fmt:
	goimports -w ./appengine

test:
	@goapp test -v -race $$(goapp list ./... | grep -v "vendor")

serve-app:
	@make -C appengine/backend serve

deploy-app-prod:
	@make -C appengine/backend deploy-prod

deploy-app-dev:
	@make -C appengine/backend deploy-dev

