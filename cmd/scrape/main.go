package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/dmmcquay/scrape"
	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Token string
}

var apiRates = flag.NewFlagSet("apirates", flag.ExitOnError)
var allCommits = flag.NewFlagSet("commits", flag.ExitOnError)
var openPRs = flag.NewFlagSet("openprs", flag.ExitOnError)
var closedPRs = flag.NewFlagSet("closedprs", flag.ExitOnError)
var top = flag.NewFlagSet("top100", flag.ExitOnError)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: scrape <command> org/repo")
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

	ro := strings.Split(os.Args[2], "/")
	if len(ro) != 2 {
		fmt.Println("poorly formated org/repo")
		return
	}
	org, repo := ro[0], ro[1]

	config := &config{}
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
		return
	}
	if missingOrg(org) || missingRepo(repo) {
		return
	}
	if allCommits.Parsed() {
		scrape.GetAllCommits(client, org, repo)
	}
	if top.Parsed() {
		scrape.Top100(client, org, repo)
	}
	if openPRs.Parsed() {
		scrape.GetPRs(client, org, repo, "open")
	}
	if closedPRs.Parsed() {
		scrape.GetPRs(client, org, repo, "closed")
	}
}

func missingRepo(repo string) bool {
	if repo == "" {
		fmt.Println("Please supply the repository using -repo option.")
		return true
	}
	return false
}

func missingOrg(org string) bool {
	if org == "" {
		fmt.Println("Please supply the orginization using -org option.")
		return true
	}
	return false
}
