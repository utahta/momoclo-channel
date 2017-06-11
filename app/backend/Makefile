serve:
	@cp config/dev.toml config/deploy.toml
	@goapp serve .

deploy-prod:
	@cp config/prod.toml config/deploy.toml
	@appcfg.py -A momoclo-channel update .

deploy-dev:
	@cp config/dev.toml config/deploy.toml
	@appcfg.py -A momoclo-channel-dev update .
