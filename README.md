# mb-web

[![GoDoc](https://godoc.org/github.com/taka-wang/mb-web?status.svg)](http://godoc.org/github.com/taka-wang/mb-web)
[![Go Report Card](https://goreportcard.com/badge/github.com/taka-wang/mb-web)](https://goreportcard.com/report/github.com/taka-wang/mb-web)
[![Docker](https://img.shields.io/badge/docker-ready-brightgreen.svg)](https://hub.docker.com/r/edgepro/mb-web)

REST web service for [psmb](https://github.com/taka-wang/psmb) written in Golang.

---

## Continuous Integration

I do continuous integration and build docker images after git push by self-hosted [drone.io](http://armdrone.cmwang.net) server for armhf platform , [circleci](http://circleci.com) server for x86 platform and [dockerhub](https://hub.docker.com/r/edgepro/mb-web) service.

| CI Server| Target    | Status                                                                                                                                                                     |
|----------|-----------|----------------------------------------------------------------------------------------------------------------------------------|
| Travis   | x86       | [![Travis](https://travis-ci.org/taka-wang/mb-web.svg)](https://travis-ci.org/taka-wang/mb-web)              |
| CircleCI | x86       | [![Circle](https://circleci.com/gh/taka-wang/mb-web.svg?style=shield)](https://circleci.com/gh/taka-wang/mb-web)               |
| Drone    | armhf     | [![Drone](http://armdrone.cmwang.net/api/badges/taka-wang/mb-web/status.svg)](http://armdrone.cmwang.net/taka-wang/mb-web)|

## Environment variables

> Why environment variable? Refer to the [12 factors](http://12factor.net/)

- CONF_WEB: config file path
- EP_BACKEND: endpoint of remote service discovery server (optional)


## Design principles

- [Separation of concerns](https://en.wikipedia.org/wiki/Separation_of_concerns)
- API-First Design
- Microservice Design
- Object-oriented Design
- 12-Factor App Design

## Golang package management

- I adopted [glide](https://glide.sh/) as package management system for this repository.

## Worker Pool Model

### Request from frontend

![uml](http://uml.cmwang.net:8000/plantuml/svg/5Sh13O0W3030LNG0QUBJRIfK848XfGrnU_LzjsRsnGAPb2Mfzd4024uNioOxRP3unagiphSAYZTk4pb2jrgWub0I2DHBU-gNuDgd--a5)

### Response from psmbtcp

![uml](http://uml.cmwang.net:8000/plantuml/svg/5Sh13O0W3030LNG0QUBJRIfK848XfGrnU_LzjsRsnGAPb2Mfzd4024uNioOxRP3unagiphSAYZTk4pb2jrgWub0I2DHBU-gNOEwN--a5)

---

## API Document

[Swagger YAML](docs/swagger.yaml)

## License

MIT
