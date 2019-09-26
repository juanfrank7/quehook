# Quehook

> Monitor open source activity with webhooks :loudspeaker:

[![Build Status](https://travis-ci.org/forstmeier/quehook.svg?branch=master)](https://travis-ci.org/forstmeier/quehook) [![Coverage Status](https://coveralls.io/repos/github/forstmeier/quehook/badge.svg?branch=master)](https://coveralls.io/github/forstmeier/quehook?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/forstmeier/quehook)](https://goreportcard.com/report/github.com/forstmeier/quehook)

## :beers: Introduction

**Quehook** provides subscribable webhooks to queries on **GitHub** data. New **[BigQuery](https://www.gharchive.org/#bigquery)** queries can be submitted by anyone via the application API endpoints and anyone can subscribe to webhook updates for any submitted queries. These queries are run against the public **[GH Archive](https://www.gharchive.org/)** datasets on an hourly basis every time the archive is updated. This allows for recurring questions regarding the open source community to be regularly answered on the newest available data.

## :octocat: Usage

For instructions on how to submit new queries or to subscribe to hourly updates, checkout the **Usage** section of the public website [here](https://forstmeier.github.io/quehook/). You can use [curl](https://curl.haxx.se/), [Postwoman](https://liyasthomas.github.io/postwoman/), or whatever other tool you want, to make the required requests to the appropriate endpoint.

## :round_pushpin: Roadmap

A simple MVP is the initial target for the launch but expanded functionality and a smoother application interface will be rolled out in the immediately subsequent versions. Below is the roadmap (although not necessary in chronological order):

- [ ] page listing all currently available queries with descriptions
- [ ] expand report target types (e.g. email, etc)

## :green_book: FAQ

Head over to the [FAQ page](https://forstmeier.github.io/quehook/faq) to get more information on the data being used and how it's being used. A changelog will be maintained on the [releases tab](https://github.com/forstmeier/quehook/releases).
