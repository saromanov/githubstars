package githubstars

import (
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strings"
)

const (
	DBNAME     = "githubstars"
	COLLECTION = "stars"
)

type GithubStars struct {
	client       *github.Client
	popularwords map[string]int
	repos        map[int]github.Repository
	currentrepos []github.Repository
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
	gs.currentrepos = []github.Repository{}
	gs.limit = 3
	return gs
}

func initMongo() *mgo.Session {
	log.Printf("Connection to the Mongodb...")
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	return sess
}

func (gs *GithubStars) Show(numstars, str, language string) {
	query := ""
	if language != "" {
		query = fmt.Sprintf("language:%s ", language)
	}
	query += fmt.Sprintf("stars:%s", numstars)
	opt := &github.SearchOptions{Sort: "stars"}
	log.Printf("Request to Github...")
	result, _, err := gs.client.Search.Repositories(query, opt)
	if err != nil {
		panic(err)
	}

	gs.currentrepos = result.Repositories
	gs.db = gs.mongosession.DB(DBNAME).C(gs.getWriteCollectionName())
	repos := []StarsInfo{}
	log.Printf("Results...")
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
		repos = append(repos, StarsInfo{*repo.FullName, *repo.StargazersCount})
	}

	gs.outputResults(repos)

}

//Commit provides write to mongodb current results
func (gs *GithubStars) Commit() {
	if len(gs.currentrepos) == 0 {
		log.Printf("Can't find current repositories for commit")
		return
	}
	db := gs.mongosession.DB(DBNAME).C(gs.getWriteCollectionName())
	db.DropCollection()
	for _, repo := range gs.currentrepos {
		gs.setData(*repo.FullName, *repo.StargazersCount)
	}

}

//THis private method provides output and comparing and formatting results
func (gs *GithubStars) outputResults(current []StarsInfo) {
	result1 := gs.getData("stars1")
	//result2 := gs.getData("stars3")

	for i, repo := range result1 {
		diff := current[i].NumStars - repo.NumStars
		diffmsg := ""
		if diff > 0 {
			diffmsg = fmt.Sprintf("(+ %d)", diff)
		} else if diff < 0 {
			diffmsg = fmt.Sprintf("(- %d", repo.NumStars-current[i].NumStars)
		}
		fmt.Println(repo.Title, repo.NumStars, current[i].NumStars, diffmsg)
	}
}

func (gs *GithubStars) getData(collname string) []StarsInfo {
	var sinfo []StarsInfo
	db := gs.mongosession.DB(DBNAME).C(gs.getWriteCollectionName())
	err := db.Find(bson.M{}).All(&sinfo)
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

//This method returns collection name for writing data
//It needs because we have limit of number of collections and
//if reading collection = limit collection, write data to
//collection1 name.I.E overwwrite data.
func (gs *GithubStars) getWriteCollectionName() string {
	return "stars1"
	/*size := gs.collectionSize()
	if size == 0 || size >= gs.limit {
		return "stars1"
	} else {
		return fmt.Sprintf("%s%d", COLLECTION, size+1)
	}*/
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
