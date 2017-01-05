# Scrape
 
# Overview

Scrape is a simple CLI tool to help Gather data from github repos about contributions

# Getting Started
If you have a working Go enviroment on your computer, you can download it by running:

`go get -u github.com/geckoboard/prism`.

# Usage

## GitHub access token

Scrape requires a GitHub access token before it can be used. See 
[Github Documentation][gat] for further information on how to get one. After you 
have the token, please set the env variable `SCRAPE_TOKEN` to the access token.

[gat]: https://help.github.com/articles/creating-an-access-token-for-command-line-use/

## scrape apirates

## Org and Repo

The following commands all require an org and Repo to be specified. An example 
would be for this repository where Org is dmmcquay and Repo is scrape.

## scrape top100

running `scrape top100 -org foo -repo bar` will return a list of the top 100 
contributors to repository.

## scrape commits

running `scrape commits -org foo -repo bar` will return a list of all contributors 
and a total count of commits for specified repository.

## scrape openprs

running `scrape openprs -org foo -repo bar` will return a list of all contributors
and a total count of open PRs they have for the specified repository.

## scrape closedprs

running `scrape closedprs -org foo -repo bar` will return a list of all contributors
and a total count of closed PRs they have for the specified repository.

