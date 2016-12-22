package main

import (
	"flag"
	"fmt"
	"log"

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

var org = flag.String("o", "", "github orginization")
var repo = flag.String("r", "", "github repository")
var check = flag.String("c", "top100", "Type of check to run (top100 (default), allcommits, apirates")

func main() {
	flag.Parse()

	if *org == "" {
		log.Fatal("need to specify an org")
	}
	if *repo == "" {
		log.Fatal("need to specify an repo")
	}

	config := &Config{}
	err := envconfig.Process("scrape", config)
	if err != nil {
		log.Fatal(err)
	}
	if config.Token == "" {
		log.Fatal("needs an access token")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	switch *check {
	case "top100":
		top100(client)
	case "allcommits":
		getAllCommits(client)
	case "apirates":
		rateLimit(client)
	case "openprs":
		getAllOpenPRs(client)
	case "closedprs":
		getAllClosedPRs(client)
	default:
		log.Fatal("not a valid check")
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

func getAllOpenPRs(client *github.Client) {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	m := make(map[string]*stat)
	for {
		prs, resp, err := client.PullRequests.List(*org, *repo, opt)
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

func getAllClosedPRs(client *github.Client) {
	opt := &github.PullRequestListOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	m := make(map[string]*stat)
	for {
		prs, resp, err := client.PullRequests.List(*org, *repo, opt)
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

func getAllCommits(client *github.Client) {
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	m := make(map[string]*stat)
	for {
		commits, resp, err := client.Repositories.ListCommits(*org, *repo, opt)
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

func top100(client *github.Client) {
	stats, _, err := client.Repositories.ListContributorsStats(*org, *repo)
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
