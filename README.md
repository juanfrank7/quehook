# Comana

> Monitor open source repository community activity :telescope:

[![Build Status](https://travis-ci.org/forstmeier/comana.svg?branch=master)](https://travis-ci.org/forstmeier/comana) [![Coverage Status](https://coveralls.io/repos/github/forstmeier/comana/badge.svg?branch=master)](https://coveralls.io/github/forstmeier/comana?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/forstmeier/comana)](https://goreportcard.com/report/github.com/forstmeier/comana)

## :beers: Introduction

**Comana** is an application to monitor open source community activity and provide data-driven insights and visibility into the projects. Data is sourced from the publicly available resources on **[GH Archive](http://www.gharchive.org/)** and processed into actionable information including raw event counts, statistics, and forward projections which can be viewed by third-parties like interested contributors or financial supporters.

## :octocat: Usage

For instructions on how to retrieve data and reports from the application, checkout the **Instructions** section of the public website [here](https://forstmeier.github.io/comana/). You can use [curl](https://curl.haxx.se/), [Postwoman](https://liyasthomas.github.io/postwoman/), or whatever other tool you want, to make a GET HTTPS request to the application API endpoint.

## :round_pushpin: Roadmap

A simple MVP is the initial target for the launch but expanded functionality and a smoother application interface will be rolled out in the immediately subsequent versions. Below is the roadmap (although not necessary in chronological order):

- [ ] embedded AWS QuickSight dashboard in landing page
  - _resources_: Athena, QuickSight, and compressed report files
  - this will also be available but likely behind a small paywall since there are costs associated with Athena
- [ ] "boosted" repo status for projects actively building communities
  - _resources_: API Gateway, DynamoDB, TTL, BuyMeACoffee, and updated API releases
  - this could be relatively easy to set up but I don't want to release it until some more datapoints/groupings are available
- [ ] expanded datapoint availability including percentages, deltas, and original statistics
  - _resources_: Step Functions, Lambda, S3, and general data analysis
  - this will be a longer-running process since it will need to include not only different types of data but different groupings of it; stuff like forward projections especially will need some outside contribution to make happen
- [ ] filtering options for fetching data from application
  -  _resources_: Step Functions, Lambda, API Gateway, updated handler logic
  - this will include the ability to filter for different amounts, groupings, and types of data; the first release will focus on the existing count datapoints

## :green_book: FAQ

Head over to the [FAQ page](https://forstmeier.github.io/comana/faq) to get more information on the data being used and how it's being used. A changelog will be maintained on the [releases tab](https://github.com/forstmeier/comana/releases).
