package backend

import
(
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"os"
	"log"
)

type Mongobackend struct {
	mongosession *mgo.Session
	db           *mgo.Collection
}

func InitMongo(addr string) *mgo.Session {
	mbackend := new(Mongobackend)
	log.Printf("Connection to the Mongodb...")
	sess, err := mgo.Dial(addr)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	mbackend.mongosession = sess
	return mbackend
}

func (mbackend *Mongobackend) SetCollection(title string) {
	mbackend.mongosession.DB(title).C(mbackend.getWriteCollectionName())
}

//This method returns collection name for writing data
func (mbackend *Mongobackend) getWriteCollectionName() string {
	return "stars1"
}

//Commit provides write to mongodb current results
//name - collection name
// If name is "", write to collection by default
func (mbackend *Mongobackend) Commit(dbname, name string) {
	/*if len(gs.currentrepos) == 0 {
		log.Fatal("Can't find current repositories for commit")
		return
	}*/

	if name == "" {
		name = mbackend.getWriteCollectionName()
	} else {

	}

	db := mbackend.mongosession.DB(dbname).C(name)
	db.DropCollection()
	mbackend.db = db
	mbackend.db.Insert(timeinfo{bson.Now()})
	for _, repo := range gs.currentrepos {
		gs.SetData(*repo.FullName, *repo.StargazersCount)
	}

	log.Printf(fmt.Sprintf("Commit new data to db %s", gs.dbname))

}

func (mbackend *Mongobackend) SetData(title string, starscount int) {
	err := mbackend.db.Insert(&starsinfo{title, starscount})
	if err != nil {
		panic(err)
	}
}