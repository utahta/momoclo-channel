# momoclo-channel

[![Build Status](https://travis-ci.org/utahta/momoclo-channel.svg?branch=master)](https://travis-ci.org/utahta/momoclo-channel)
[![Go Report Card](https://goreportcard.com/badge/github.com/utahta/momoclo-channel)](https://goreportcard.com/report/github.com/utahta/momoclo-channel)

This is an app that gives you any information about Momoiro Clover Z with LINE and Twitter.

- [LINE](https://momoclo-channel.com/line/bot/about)
- [@botnofu](https://twitter.com/botnofu)

## Architecture

To respect the [Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html)

| Layer | Directory |
| --- | --- |
| External interfaces | [infrastructure](./infrastructure) |
| Interface adapters / Controllers | [adapter/handler](./adapter/handler) |
| Interface adapters / Gateways | [adapter/gateway](./adapter/gateway) |
| Use Cases | [usecase](./usecase) |
| Entities | [domain](./domain) |
