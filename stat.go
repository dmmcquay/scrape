package scrape

import "fmt"

type stat struct {
	Login string   `json:"login"`
	Email []string `json:"email"`
	Count int      `json:"count"`
	Rank  int      `json:"rank"`
}

type byCount []stat

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

func (s byCount) Len() int {
	return len(s)
}

func (s byCount) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byCount) Less(i, j int) bool {
	return s[i].Count < s[j].Count
}
