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
	//DBNAME current db in mongodb
	DBNAME = "githubstars"
	//COLLECTION should be replaced to search title
	COLLECTION = "stars1"
)

type githubstars struct {
	client       *github.Client
	popularwords map[string]int
	repos        map[int]github.Repository
	currentrepos []github.Repository
	mongosession *mgo.Session
	db           *mgo.Collection
	limit        int
	dbname       string
}

type StarsInfo struct {
	Title    string
	NumStars int
}

//summary provides some stat information before output
type summary struct {
	most         record
	fewest_stars record
}

type record struct {
	title string
	item  int
}

func Init() *githubstars {
	gs := new(githubstars)
	gs.client = github.NewClient(nil)
	gs.popularwords = map[string]int{}
	gs.repos = map[int]github.Repository{}
	gs.mongosession = initMongo()
	gs.currentrepos = []github.Repository{}
	gs.limit = 3
	gs.dbname = ""
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

//Show provides output information
func (gs *githubstars) Show(numstars, str, language string) {
	query := ""
	dbname := ""
	if language != "" {
		query = fmt.Sprintf("language:%s ", language)
		dbname += language
	}
	query += fmt.Sprintf("stars:%s", numstars)
	dbname += str
	dbname += numstars
	gs.dbname = constructName(dbname)
	opt := &github.SearchOptions{Sort: "stars"}
	log.Printf("Request to Github...")
	result, _, err := gs.client.Search.Repositories(query, opt)
	if err != nil {
		panic(err)
	}

	gs.currentrepos = result.Repositories
	gs.db = gs.mongosession.DB(gs.dbname).C(gs.getWriteCollectionName())
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
	gs.outputResults(repos, gs.dbname, COLLECTION)

}

//Commit provides write to mongodb current results
func (gs *githubstars) Commit() {
	if len(gs.currentrepos) == 0 {
		log.Fatal("Can't find current repositories for commit")
		return
	}
	db := gs.mongosession.DB(gs.dbname).C(gs.getWriteCollectionName())
	db.DropCollection()
	for _, repo := range gs.currentrepos {
		gs.setData(*repo.FullName, *repo.StargazersCount)
	}

}

//CompareWith provides comparation with results
func (gs *githubstars) CompareWith(dbtitle string) {
	data := gs.getData(dbtitle, COLLECTION)
	if len(data) == 0 {
		log.Fatal(fmt.Sprintf("Collection %s is not found", dbtitle))
	}
	alldata := []StarsInfo{}
	for _, value := range gs.repos {
		alldata = append(alldata, StarsInfo{*value.FullName, *value.StargazersCount})
	}
	gs.outputResults(alldata, dbtitle, COLLECTION)

}

//This private method provides output and comparing and formatting results
func (gs *githubstars) outputResults(current []StarsInfo, dbname string, collname string) {
	result1 := gs.getData(dbname, collname)
	if len(result1) == 0 {
		log.Printf(fmt.Sprintf("db %s or collection %s not found", dbname, collname))
	}
	//result2 := gs.getData("stars3")
	summ := summary{}
	summ.most = record{}
	summ.fewest_stars = record{}
	summ.fewest_stars.item = 99999999
	for i, repo := range result1 {
		diff := current[i].NumStars - repo.NumStars
		diffmsg := ""
		if summ.most.item < diff {
			summ.most.item = diff
			summ.most.title = repo.Title
		}

		if summ.fewest_stars.item >= diff {
			summ.fewest_stars.item = diff
			summ.fewest_stars.title = repo.Title
		}

		if diff > 0 {
			diffmsg = fmt.Sprintf("(+ %d)", diff)
		} else if diff < 0 {
			diffmsg = fmt.Sprintf("(- %d", repo.NumStars-current[i].NumStars)
		}
		fmt.Println(repo.Title, repo.NumStars, current[i].NumStars, diffmsg)
	}

	log.Printf("Summary...")
	fmt.Println(" ")
	fmt.Println(fmt.Sprintf("Most number of new stars: %s %d", summ.most.title, summ.most.item))
	fmt.Println(fmt.Sprintf("Fewest number of new stars: %s %d",
		summ.fewest_stars.title, summ.fewest_stars.item))
}

//get data from mongo
func (gs *githubstars) getData(dbname, collname string) []StarsInfo {
	var sinfo []StarsInfo
	db := gs.mongosession.DB(dbname).C(collname)
	err := db.Find(bson.M{}).All(&sinfo)
	if err != nil {
		panic(err)
	}
	return sinfo
}

//return number of records from collection
func (gs *githubstars) collectionSize() int {
	count, err := gs.mongosession.DB(DBNAME).CollectionNames()
	if err != nil {
		return 0
	}
	return len(count)
}

//This method returns collection name for writing data
func (gs *githubstars) getWriteCollectionName() string {
	return "stars1"
}

func (gs *githubstars) setData(title string, starscount int) {
	err := gs.db.Insert(&StarsInfo{title, starscount})
	if err != nil {
		panic(err)
	}

}

func splitDescription(desc string) []string {
	result := strings.Split(desc, " ")
	return result
}

//This method provides construction valid name for mongodb
//For example: if we have query ">1000",
//constructName will output "gr1000"
func constructName(title string) string {
	if strings.Index(title, ">") != -1 {
		return strings.Replace(title, ">", "gr", 1)
	}

	if strings.Index(title, "<") != -1 {
		return strings.Replace(title, "<", "lo", 1)
	}

	return title
}
