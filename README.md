# Quehook

> Monitor open source activity with webhooks :loudspeaker:

[![Build Status](https://travis-ci.org/forstmeier/quehook.svg?branch=master)](https://travis-ci.org/forstmeier/quehook) [![Coverage Status](https://coveralls.io/repos/github/forstmeier/quehook/badge.svg?branch=master)](https://coveralls.io/github/forstmeier/quehook?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/forstmeier/quehook)](https://goreportcard.com/report/github.com/forstmeier/quehook)

<style>.bmc-button img{width: 27px !important;margin-bottom: 1px !important;box-shadow: none !important;border: none !important;vertical-align: middle !important;}.bmc-button{line-height: 36px !important;height:37px !important;text-decoration: none !important;display:inline-flex !important;color:#FFFFFF !important;background-color:#FF813F !important;border-radius: 3px !important;border: 1px solid transparent !important;padding: 1px 9px !important;font-size: 22px !important;letter-spacing: 0.6px !important;box-shadow: 0px 1px 2px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 1px 2px 2px rgba(190, 190, 190, 0.5) !important;margin: 0 auto !important;font-family:'Cookie', cursive !important;-webkit-box-sizing: border-box !important;box-sizing: border-box !important;-o-transition: 0.3s all linear !important;-webkit-transition: 0.3s all linear !important;-moz-transition: 0.3s all linear !important;-ms-transition: 0.3s all linear !important;transition: 0.3s all linear !important;}.bmc-button:hover, .bmc-button:active, .bmc-button:focus {-webkit-box-shadow: 0px 1px 2px 2px rgba(190, 190, 190, 0.5) !important;text-decoration: none !important;box-shadow: 0px 1px 2px 2px rgba(190, 190, 190, 0.5) !important;opacity: 0.85 !important;color:#FFFFFF !important;}</style><link href="https://fonts.googleapis.com/css?family=Cookie" rel="stylesheet"><a class="bmc-button" target="_blank" href="https://www.buymeacoffee.com/forstmeier"><img src="https://bmc-cdn.nyc3.digitaloceanspaces.com/BMC-button-images/BMC-btn-logo.svg" alt="Buy me a coffee"><span style="margin-left:5px">Buy me a coffee</span></a>

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
