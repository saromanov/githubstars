package githubstars

import (
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strings"
)

const (
	DBNAME = "githubstars"
)

type GithubStars struct {
	client       *github.Client
	popularwords map[string]int
	repos        map[int]github.Repository
	mongosession *mgo.Session
	db           *mgo.Collection
	limit        int
}

type StarsInfo struct {
	Title    string
	NumStars int
}

func Init() *GithubStars {
	gs := new(GithubStars)
	gs.client = github.NewClient(nil)
	gs.popularwords = map[string]int{}
	gs.repos = map[int]github.Repository{}
	gs.mongosession = initMongo()
	gs.db = gs.mongosession.DB(DBNAME).C("stars1")
	gs.limit = 3
	return gs
}

func initMongo() *mgo.Session {
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	return sess
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
		gs.setData(*repo.FullName, *repo.StargazersCount)
	}

	gs.getData()
}

func (gs *GithubStars) getData() []StarsInfo {
	var sinfo []StarsInfo
	err := gs.db.Find(bson.M{}).All(&sinfo)
	if err != nil {
		panic(err)
	}
	return sinfo
}

func (gs *GithubStars) collectionSize() int {
	count, err := gs.mongosession.DB(DBNAME).CollectionNames()
	if err != nil {
		return 0
	}
	return len(count)
}

func (gs *GithubStars) setData(title string, starscount int) {
	err := gs.db.Insert(&StarsInfo{title, starscount})
	if err != nil {
		panic(err)
	}

}

func splitDescription(desc string) []string {
	result := strings.Split(desc, " ")
	return result
}
