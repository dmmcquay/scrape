package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"golang.org/x/oauth2"

	"github.com/dmmcquay/scrape"
	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
)

type stat struct {
	Login string   `json:"login"`
	Email []string `json:"email"`
	Count int      `json:"count"`
	Rank  int      `json:"rank"`
}

type byCount []stat

func (s byCount) Len() int           { return len(s) }
func (s byCount) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byCount) Less(i, j int) bool { return s[i].Count < s[j].Count }

type Config struct {
	Token string
}

func (s stat) String() string {
	if len(s.Email) < 3 {
		return fmt.Sprintf("%d\t%s\t%v\t%d", s.Rank, s.Login, s.Email, s.Count)
	}
	return fmt.Sprintf(
		"%d\t%s\t%v\t%d",
		s.Rank,
		s.Login,
		fmt.Sprintf("[%s [...] %s]", s.Email[0], s.Email[len(s.Email)-1]),
		s.Count,
	)
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
		getPRs(client, *openPRsOrg, *openPRsRepo, "open")
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
		getPRs(client, *closedPRsOrg, *closedPRsRepo, "closed")
	}
}

func checkAndAddEmail(e string, emails []string) bool {
	for _, s := range emails {
		if e == s {
			return true
		}
	}
	return false
}

func getPRs(client *github.Client, org, repo, state string) {
	opt := &github.PullRequestListOptions{
		State: state,
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
				m[a] = &stat{Login: a, Email: []string{}, Count: 1}
				continue
			}
			tmp := m[a]
			tmp.Count += 1
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	b := byCount{}
	for _, v := range m {
		b = append(b, *v)
	}
	sort.Sort(b)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 10, 8, 0, '\t', 0)
	fmt.Fprintln(w, "rank\tlogin\tPRs")

	total := 0
	atotal := len(m)
	for n, v := range b {
		total += v.Count
		v.Rank = atotal - n
		fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%d", v.Rank, v.Login, v.Count))
	}
	fmt.Fprintln(w)
	w.Flush()
	fmt.Printf("TOTAL PRs: %d\n", total)
	fmt.Printf("TOTAL AUTHORS: %d\n", atotal)
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
				m[a] = &stat{Login: a, Email: []string{e}, Count: 1}
				continue
			}
			tmp := m[a]
			tmp.Count += 1
			if !checkAndAddEmail(e, tmp.Email) {
				tmp.Email = append(tmp.Email, e)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	b := byCount{}
	for _, v := range m {
		b = append(b, *v)
	}
	sort.Sort(b)

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 10, 8, 0, '\t', 0)
	fmt.Fprintln(w, "rank\tlogin\temails\tcommits")

	total := 0
	atotal := len(m)
	for n, v := range b {
		total += v.Count
		v.Rank = atotal - n
		fmt.Fprintln(w, v)
	}
	fmt.Fprintln(w)
	w.Flush()
	fmt.Printf("TOTAL COMMITS: %d\n", total)
	fmt.Printf("TOTAL AUTHORS: %d\n", atotal)
}
