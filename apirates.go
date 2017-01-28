package scrape

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
)

func RateLimit(client *github.Client) {
	r, _, err := client.RateLimit()
	if err != nil {
		log.Printf("error getting rate: %v", err)
		return
	}
	fmt.Printf("%d/%d requests\n", r.Remaining, r.Limit)
}
