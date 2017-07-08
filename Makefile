install:
	@dep ensure

fmt:
	@goimports -w $$(goapp list -f '{{.Dir}}' ./... | grep -v "vendor")

test:
	@goapp test -v -race $$(goapp list ./... | grep -v "vendor")

serve:
	@make -C app/backend prepare-serve
	@goapp serve ./app/backend/app.yaml ./app/queue/app.yaml

deploy-prod:
	@make -C app/queue deploy-prod
	@make -C app/backend deploy-prod

deploy-dev:
	@make -C app/queue deploy-dev
	@make -C app/backend deploy-dev
