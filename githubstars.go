package githubstars

import (
	"fmt"
	"github.com/google/go-github/github"
	"strings"
)

type GithubStars struct {
	client       *github.Client
	popularwords map[string]int
	repos        map[int]github.Repository
}

func Init() *GithubStars {
	gs := new(GithubStars)
	gs.client = github.NewClient(nil)
	gs.popularwords = map[string]int{}
	gs.repos = map[int]github.Repository{}
	return gs
}

func (gs *GithubStars) Get(numstars, str, language string) {
	query := ""
	if language != "" {
		query = fmt.Sprintf("language:%s ", language)
	}
	query += fmt.Sprintf("stars:%s", numstars)
	opt := &github.SearchOptions{Sort: "stars"}
	result, _, err := gs.client.Search.Repositories(query, opt)
	if err != nil {
		panic(err)
	}

	for i, repo := range result.Repositories {
		words := splitDescription(*repo.Description)
		for _, word := range words {
			_, ok := gs.popularwords[word]
			if !ok {
				gs.popularwords[word] = 0
			} else {
				gs.popularwords[word] += 1
			}
		}
		gs.repos[i] = repo
		//fmt.Println(*repo.FullName, *repo.StargazersCount, splitDescription(*repo.Description))
	}
}

func splitDescription(desc string) []string {
	result := strings.Split(desc, " ")
	return result
}
