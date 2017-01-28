package scrape

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/google/go-github/github"
)

// Top100 prints to stdout the top100 contributors to organization's
// repository
func Top100(client *github.Client, org, repo string) {
	stats, _, err := client.Repositories.ListContributorsStats(org, repo)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Println("hit rate limit")
	}
	if err != nil {
		log.Fatal(err)
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 10, 8, 0, '\t', 0)
	fmt.Fprintln(w, "rank\tlogin\tcommits")

	for n, i := range stats {
		fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%d", (100-n), *i.Author.Login, *i.Total))
	}
	fmt.Fprintln(w)
	w.Flush()
	fmt.Printf("TOTAL TOP100 AUTHORS: %d\n", len(stats))
}
