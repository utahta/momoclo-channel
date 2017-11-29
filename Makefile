install:
	@dep ensure

fmt:
	@goimports -w $$(goapp list -f '{{.Dir}}' ./... | grep -v "vendor")

test:
	@goapp test -v $$(goapp list ./... | grep -v "vendor")

lint:
	@golint $$(go list ./... | grep -v vendor) | grep -v ": exported const" | grep -v ": exported var Err"

review: test
	@make lint | reviewdog -f=golint -diff="git diff master"

serve:
	@make -C appengine/backend prepare-serve
	@goapp serve ./appengine/backend/app.yaml ./appengine/queue/app.yaml

deploy-prod:
	@make -C appengine/batch deploy-prod
	@make -C appengine/backend deploy-prod

deploy-dev:
	@make -C appengine/batch deploy-dev
	@make -C appengine/backend deploy-dev
