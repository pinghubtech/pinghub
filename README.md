
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/traefik/traefik/blob/master/LICENSE.md)
![v0.0.1](https://github.com/antgubarev/pingbot/actions/workflows/quality.yml/badge.svg?branch=v0.0.1)
![Latest release](https://img.shields.io/github/v/release/antgubarev/pingbot)
[![Go Report Card](https://goreportcard.com/badge/github.com/antgubarev/pingbot)](https://goreportcard.com/report/github.com/antgubarev/pingbot)

#About

PingHub is open-source uptime monitoring which ping your site many times in every minute, 
fixes downtime and calculate avg response time.

# Architecture

PingHub include two parts:
 - checker that checks targets status and stores results
 - http server that renders status page

You have to run both part separate applications

# Installation

**Binary**

Grab the latest binary from the [releases page](https://github.com/antgubarev/pinghub/releases)

**Docker**

See example `docker-compose.dev.yml`

#Configuration

You may set environment variables to configure PingBot

- `HTTP_LISTEN_ADDR` (default `localhost:8080`). Public addr for status pages
- `LOG_LEVEL` (default `error`). Possible values: `debug`, `info`, `warn`, `error`, `fatal`, `panic`

# Contributing

- Fork it
- Create your feature branch (git checkout -b my-new-feature)
- Commit your changes (git commit -am 'Added some feature')
- Push to the branch (git push origin my-new-feature)
- Create new Pull Request
