package scrape

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/github"
)

// RateLimit prints to stdout number of api request remaining
func RateLimit(client *github.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, _, err := client.RateLimits(ctx)
	if err != nil {
		log.Printf("error getting rate: %v", err)
		return
	}
	fmt.Printf("%d/%d requests\n", r.Core.Remaining, r.Core.Limit)
}
