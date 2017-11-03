install:
	@dep ensure

fmt:
	@goimports -w $$(goapp list -f '{{.Dir}}' ./... | grep -v "vendor")

test:
	@go test -v -race $$(goapp list ./... | grep -v "vendor")

serve:
	@make -C appengine/backend prepare-serve
	@goapp serve ./appengine/backend/app.yaml ./appengine/queue/app.yaml

deploy-prod:
	@make -C appengine/queue deploy-prod
	@make -C appengine/backend deploy-prod

deploy-dev:
	@make -C appengine/queue deploy-dev
	@make -C appengine/backend deploy-dev
