# mb-web

[![Docker](https://img.shields.io/badge/docker-ready-brightgreen.svg)](https://hub.docker.com/r/edgepro/mb-web)
[![GoDoc](https://godoc.org/github.com/taka-wang/mb-web?status.svg)](http://godoc.org/github.com/taka-wang/mb-web)
[![Go Report Card](https://goreportcard.com/badge/github.com/taka-wang/mb-web)](https://goreportcard.com/report/github.com/taka-wang/mb-web)

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

- [Separation of concerns](https://en.wikipedia.org/wiki/Separation_of_concerns) - I separate config, route and worker functions to respective packages, link route and worker services by predefined function signature and go channel.
- [API-First Design](http://www.api-first.com/)
- [Microservice Design](https://en.wikipedia.org/wiki/Microservices)
- [Object-oriented Design](https://en.wikipedia.org/wiki/Object-oriented_design)
- [12-Factor App Design](http://12factor.net/)

## REST API Testing

I implement **three** kinds of REST API testing.

- [Go testing](test/client_test.go)
- [Shell script tester](image/shell.gif)    
- [Postman/Newman tester](image/newman.gif)

![gif](image/newman.gif)

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
