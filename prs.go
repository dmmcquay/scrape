package scrape

import (
	"fmt"
	"log"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/google/go-github/github"
)

// GetPRs prints to stdout a sorted list of either closed or open PRs to
// specified organization's repository
func GetPRs(client *github.Client, org, repo, state string) {
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
