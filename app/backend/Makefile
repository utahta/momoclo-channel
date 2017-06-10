serve:
	@cp .env.local env
	@goapp serve .

deploy-prod:
	@cp .env.prod env
	@appcfg.py -A momoclo-channel update .

deploy-dev:
	@cp .env.dev env
	@appcfg.py -A momoclo-channel-dev update .
