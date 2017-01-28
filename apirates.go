package scrape

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
)

// RateLimit prints to stdout number of api request remaining
func RateLimit(client *github.Client) {
	r, _, err := client.RateLimit()
	if err != nil {
		log.Printf("error getting rate: %v", err)
		return
	}
	fmt.Printf("%d/%d requests\n", r.Remaining, r.Limit)
}
