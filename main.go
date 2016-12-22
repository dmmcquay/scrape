package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
)

type stat struct {
	email []string
	count int
}

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
		rateLimit(client)
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
		getAllCommits(client, *allCommitsOrg, *allCommitsRepo)
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
		top100(client, *topOrg, *topRepo)
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
		getAllOpenPRs(client, *openPRsOrg, *openPRsRepo)
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
		getAllClosedPRs(client, *closedPRsOrg, *closedPRsRepo)
	}
}

func rateLimit(client *github.Client) {
	r, _, err := client.RateLimit()
	if err != nil {
		log.Printf("error getting rate: %v", err)
		return
	}
	fmt.Printf("%d/%d\n", r.Remaining, r.Limit)
}

func checkAndAddEmail(e string, emails []string) bool {
	for _, s := range emails {
		if e == s {
			return true
		}
	}
	return false
}

func getAllOpenPRs(client *github.Client, org, repo string) {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	m := make(map[string]*stat)
	for {
		prs, resp, err := client.PullRequests.List(org, repo, opt)
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, pr := range prs {
			a := "username missing"
			if pr.User != nil {
				a = *pr.User.Login
			}
			_, ok := m[a]
			if !ok {
				m[a] = &stat{email: []string{}, count: 1}
				continue
			}
			tmp := m[a]
			tmp.count += 1
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	total := 0
	atotal := 0
	for k, v := range m {
		total += v.count
		atotal += 1
		fmt.Printf("Author: %s count: %d\n", k, v.count)
	}
	fmt.Printf("TOTAL: %d\n", total)
	fmt.Printf("ATOTAL: %d\n", atotal)
}

func getAllClosedPRs(client *github.Client, org, repo string) {
	opt := &github.PullRequestListOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	m := make(map[string]*stat)
	for {
		prs, resp, err := client.PullRequests.List(org, repo, opt)
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, pr := range prs {
			a := "username missing"
			if pr.User != nil {
				a = *pr.User.Login
			}
			_, ok := m[a]
			if !ok {
				m[a] = &stat{email: []string{}, count: 1}
				continue
			}
			tmp := m[a]
			tmp.count += 1
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	total := 0
	atotal := 0
	for k, v := range m {
		total += v.count
		atotal += 1
		fmt.Printf("Author: %s count: %d\n", k, v.count)
	}
	fmt.Printf("TOTAL: %d\n", total)
	fmt.Printf("ATOTAL: %d\n", atotal)
}

func getAllCommits(client *github.Client, org, repo string) {
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	m := make(map[string]*stat)
	for {
		commits, resp, err := client.Repositories.ListCommits(org, repo, opt)
		if _, ok := err.(*github.RateLimitError); ok {
			log.Println("hit rate limit")
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range commits {
			a := "username missing"
			e := "fake@fake.com"
			if c.Author != nil {
				a = *c.Author.Login
			}
			if c.Commit.Author != nil {
				e = *c.Commit.Author.Email
			}
			_, ok := m[a]
			if !ok {
				m[a] = &stat{email: []string{e}, count: 1}
				continue
			}
			tmp := m[a]
			tmp.count += 1
			if !checkAndAddEmail(e, tmp.email) {
				tmp.email = append(tmp.email, e)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	total := 0
	atotal := 0
	for k, v := range m {
		total += v.count
		atotal += 1
		fmt.Printf("Author: %s count: %v\n", k, v)
	}
	fmt.Printf("TOTAL: %d\n", total)
	fmt.Printf("ATOTAL: %d\n", atotal)
}

func top100(client *github.Client, org, repo string) {
	stats, _, err := client.Repositories.ListContributorsStats(org, repo)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("hit rate limit")
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(stats))
	for _, i := range stats {
		fmt.Printf(
			"Author: %v commits: %v\n",
			*i.Author.Login,
			*i.Total,
		)
	}
}
