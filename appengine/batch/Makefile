prepare-serve:
	@cp config/dev.toml config/deploy.toml

serve: prepare-serve
	@goapp serve .

deploy-prod:
	@cp config/prod.toml config/deploy.toml
	@appcfg.py -A momoclo-channel update .

deploy-dev:
	@cp config/dev.toml config/deploy.toml
	@appcfg.py -A momoclo-channel-dev update .
