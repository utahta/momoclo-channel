install:
	@dep ensure

fmt:
	goimports -w $$(goapp list ./... | grep -v "vendor")

test:
	@goapp test -v -race $$(goapp list ./... | grep -v "vendor")

serve-app:
	@make -C app/backend serve

deploy-app-prod:
	@make -C app/backend deploy-prod

deploy-app-dev:
	@make -C app/backend deploy-dev

