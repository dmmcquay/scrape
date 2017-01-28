package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"

	"github.com/dmmcquay/scrape"
	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Token string
}

var apiRates = flag.NewFlagSet("apirates", flag.ExitOnError)

var allCommits = flag.NewFlagSet("commits", flag.ExitOnError)
var allCommitsOrg = allCommits.String("org", "", "github orginization")
var allCommitsRepo = allCommits.String("repo", "", "github repository")

var openPRs = flag.NewFlagSet("openprs", flag.ExitOnError)
var openPRsOrg = openPRs.String("org", "", "github orginization")
var openPRsRepo = openPRs.String("repo", "", "github repository")

var closedPRs = flag.NewFlagSet("closedprs", flag.ExitOnError)
var closedPRsOrg = closedPRs.String("org", "", "github orginization")
var closedPRsRepo = closedPRs.String("repo", "", "github repository")

var top = flag.NewFlagSet("top100", flag.ExitOnError)
var topOrg = top.String("org", "", "github orginization")
var topRepo = top.String("repo", "", "github repository")

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: scrape <command> [<args>]")
		fmt.Println("The scrape commands are: ")
		fmt.Println(" top100     See top 100 commiters to project")
		fmt.Println(" commits    See all user's commits to project")
		fmt.Println(" apirates   See current used api requests/total")
		fmt.Println(" openprs    See all open PRs to project")
		fmt.Println(" closedprs  See all closed PRs to project")
		return
	}

	switch os.Args[1] {
	case "apirates":
		apiRates.Parse(os.Args[2:])
	case "commits":
		allCommits.Parse(os.Args[2:])
	case "openprs":
		openPRs.Parse(os.Args[2:])
	case "closedprs":
		closedPRs.Parse(os.Args[2:])
	case "top100":
		top.Parse(os.Args[2:])
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	config := &Config{}
	err := envconfig.Process("scrape", config)
	if err != nil {
		log.Fatal(err)
	}
	if config.Token == "" {
		fmt.Println("scrape requires SCRAPE_TOKEN env variable to be defined with valid access token")
		os.Exit(3)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	if apiRates.Parsed() {
		scrape.RateLimit(client)
	}
	if allCommits.Parsed() {
		if *allCommitsOrg == "" {
			fmt.Println("Please supply the orginization using -org option.")
			return
		}

		if *allCommitsRepo == "" {
			fmt.Println("Please supply the repository using -repo option.")
			return
		}
		scrape.GetAllCommits(client, *allCommitsOrg, *allCommitsRepo)
	}
	if top.Parsed() {
		if *topOrg == "" {
			fmt.Println("Please supply the orginization using -org option.")
			return
		}

		if *topRepo == "" {
			fmt.Println("Please supply the repository using -repo option.")
			return
		}
		scrape.Top100(client, *topOrg, *topRepo)
	}
	if openPRs.Parsed() {
		if *openPRsOrg == "" {
			fmt.Println("Please supply the orginization using -org option.")
			return
		}

		if *openPRsRepo == "" {
			fmt.Println("Please supply the repository using -repo option.")
			return
		}
		scrape.GetPRs(client, *openPRsOrg, *openPRsRepo, "open")
	}
	if closedPRs.Parsed() {
		if *closedPRsOrg == "" {
			fmt.Println("Please supply the orginization using -org option.")
			return
		}

		if *closedPRsRepo == "" {
			fmt.Println("Please supply the repository using -repo option.")
			return
		}
		scrape.GetPRs(client, *closedPRsOrg, *closedPRsRepo, "closed")
	}
}
