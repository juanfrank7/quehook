# Comana

> Monitor open source repository community activity :telescope:

[![Build Status](https://travis-ci.org/forstmeier/comana.svg?branch=master)](https://travis-ci.org/forstmeier/comana) [![Coverage Status](https://coveralls.io/repos/github/forstmeier/comana/badge.svg?branch=master)](https://coveralls.io/github/forstmeier/comana?branch=master) 

## Introduction :beers:

**Comana** is an application to monitor open source repository community activity and provide data-driven insights and visibility into subscribed projects. Project maintainers can subscribe their repositories to the application which will actively collect data and make it available for third parties like interested contributors or financial supporters.

## Usage :octocat:

> The app is currently under active development but once the MVP is ready, this section will be updated :construction:

## Roadmap :round_pushpin:

A simple MVP is the initial target for the launch but expanded functionality and a smoother application interface will be rolled out in the immediately subsequent versions. Below is the roadmap (although not necessary in chronological order):

- [ ] generalized backfill data creation service
  - [ ] resources: Lambda, S3, and refactored service logic
- [ ] embedded AWS QuickSight dashboard in landing page
  - [ ] resources: Athena, QuickSight, and compressed report files
- [ ] "boosted" repo status for projects actively building communities
  - [ ] resources: API Gateway, DynamoDB, TTL, BuyMeACoffee, and updated API releases
- [ ] expanded datapoint availability including percentages, deltas, and original statistics
  - [ ] resources: Step Functions, Lambda, S3, and general data analysis
- [ ] filtering options for fetching data from application
  - [ ] resources: Step Functions, Lambda, API Gateway, updated handler logic

## FAQ :book:

> This section will eventually contain a link to the FAQ page on the application :construction:
