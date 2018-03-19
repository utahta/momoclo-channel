## setup

setup:
	@dep ensure -v

setup/tools:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u honnef.co/go/tools/cmd/unused
	go get -u github.com/kisielk/errcheck
	go get -u golang.org/x/lint/golint
	go get -u github.com/haya14busa/reviewdog/cmd/reviewdog

## test

TESTPKGS=$(shell go list ./... | grep -v "vendor")

test:
	@goapp test -v $(TESTPKGS)

## lint

LINTPKGS=$(shell go list ./... | grep -v "vendor")

lint: lint/vet

lint/vet:
	@go vet $(LINTPKGS)

## reviewdog

reviewdog:
	reviewdog -diff="git diff master"

reviewdog/ci:
	reviewdog -ci="travis"

## serve

serve:
	@make -C appengine/backend prepare-serve
	@goapp serve ./appengine/backend/app.yaml ./appengine/batch/app.yaml

## deploy

deploy-prod:
	@make -C appengine/batch deploy-prod
	@make -C appengine/backend deploy-prod

deploy-dev:
	@make -C appengine/batch deploy-dev
	@make -C appengine/backend deploy-dev
