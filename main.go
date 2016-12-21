package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

type commits []struct {
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
}

type contribs []struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}

type stat struct {
	email   []string
	commits int
}

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "27e8ffceaaaab8310e5564bea7e8a028bc181f3e"},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	//top100()
	getAllCommits(client)
	//rateLimit(client)
	//m := make(map[string]int)
	//for _, i := range listContributors() {
	//	//fmt.Println(i.Commit.Author.Name)
	//	a := i.Login
	//	_, ok := m[a]
	//	if !ok {
	//		m[a] = i.Contributions
	//		continue
	//	}
	//	m[a] += 1
	//}
	//fmt.Println(m)
	//top100()
}

func checkAndAddEmail(e string, emails []string) bool {
	for _, s := range emails {
		if e == s {
			return true
		}
	}
	return false
}

func getAllCommits(client *github.Client) {
	opt := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	//var allCommits []*github.RepositoryCommit
	m := make(map[string]*stat)
	for {
		commits, resp, err := client.Repositories.ListCommits("dmmcquay", "sqrl", opt)
		//commits, resp, err := client.Repositories.ListCommits("kubernetes", "kubernetes", opt)
		//commits, resp, err := client.Repositories.ListCommits("coreos", "dbtester", opt)
		//commits, resp, err := client.Repositories.ListCommits("coreos", "dex", opt)
		//commits, resp, err := client.Repositories.ListCommits("coreos", "torus", opt)
		//commits, resp, err := client.Repositories.ListCommits("coreos", "rkt", opt)
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
				m[a] = &stat{email: []string{e}, commits: 1}
				continue
			}
			tmp := m[a]
			tmp.commits += 1
			if !checkAndAddEmail(e, tmp.email) {
				tmp.email = append(tmp.email, e)
			}
		}
		//allCommits = append(allCommits, commits...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	total := 0
	atotal := 0
	for k, v := range m {
		total += v.commits
		atotal += 1
		fmt.Printf("Author: %s Commits: %v\n", k, v)
	}
	fmt.Printf("TOTAL: %d\n", total)
	fmt.Printf("ATOTAL: %d\n", atotal)
}

func listContributors() contribs {
	r, err := http.Get("https://api.github.com/repos/kubernetes/kubernetes/contributors")
	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode != 200 {
		log.Fatal("Unexpected status code", r.StatusCode)
	}
	defer r.Body.Close()

	target := contribs{}

	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		log.Fatal("error parsing response: %v", err)
	}
	return target
}

func listCommits() commits {
	r, err := http.Get("https://api.github.com/repos/kubernetes/kubernetes/commits")
	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode != 200 {
		log.Fatal("Unexpected status code", r.StatusCode)
	}
	defer r.Body.Close()

	target := commits{}

	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		log.Fatal("error parsing response: %v", err)
	}
	return target
}

func rateLimit(client *github.Client) {
	r, _, err := client.RateLimit()
	if err != nil {
		log.Printf("error getting rate: %v", err)
		return
	}
	fmt.Printf("%d/%d\n", r.Remaining, r.Limit)
}

func top100(client *github.Client) {
	//stats, _, err := client.Repositories.ListContributorsStats("dmmcquay", "sqrl")
	stats, _, err := client.Repositories.ListContributorsStats("kubernetes", "kubernetes")
	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("hit rate limit")
	}
	if err != nil {
		log.Fatal(err)
	}
	//if _, ok := err.(*github.AcceptedError); ok {
	//	log.Println("scheduled on GitHub side")
	//}
	fmt.Println(len(stats))
	for _, i := range stats {
		fmt.Printf(
			"Author: %v commits: %v\n",
			*i.Author.Login,
			*i.Total,
		)
	}
}
