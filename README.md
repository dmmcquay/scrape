# Scrape

[![GoDoc](https://godoc.org/github.com/dmmcquay/scrape?status.svg)](https://godoc.org/github.com/dmmcquay/scrape)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmmcquay/scrape)](https://goreportcard.com/report/github.com/dmmcquay/scrape)
 
# Overview

Scrape is a simple CLI tool to help Gather data from github repos about contributions

# Getting Started
If you have a working Go enviroment on your computer, you can download it by running:

`go get -u github.com/dmmcquay/scrape`

# Usage

## GitHub access token

Scrape requires a GitHub access token before it can be used. See 
[Github Documentation][gat] for further information on how to get one. After you 
have the token, please set the env variable `SCRAPE_TOKEN` to the access token.

[gat]: https://help.github.com/articles/creating-an-access-token-for-command-line-use/

## scrape apirates

## Org and Repo

The following commands all require an org and Repo to be specified. An example 
would be for this repository where Org is dmmcquay and Repo is scrape. The format
for this would be `dmmcquay/scrape` 

## scrape top100

running: 

```
scrape top100 foo/bar
``` 

will return a list of the top 100 contributors to repository.

## scrape commits

running:

```
scrape commits foo/bar
```

will return a list of all contributors and a total count of commits for 
specified repository.

## scrape openprs

running:

```
scrape openprs foo/bar
``` 

will return a list of all contributors and a total count of open PRs they have 
for the specified repository.

## scrape closedprs

running:

```
scrape closedprs foo/bar
``` 
will return a list of all contributors and a total count of closed PRs they 
have for the specified repository.

