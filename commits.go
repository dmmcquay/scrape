package scrape

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/google/go-github/github"
)

func checkAndAddEmail(e string, emails []string) bool {
	for _, s := range emails {
		if e == s {
			return true
		}
	}
	return false
}

// GetAllCommits prints to stdout a sorted list of all commits to a
// specified organization's repository
func GetAllCommits(client *github.Client, org, repo string) {
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
