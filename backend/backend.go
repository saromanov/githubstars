package backend

//Backend provides interface for backends for storing info
type Backend interface{
	//SetCollection for choose target collection
	SetCollection(title string)

	//Commit data
	Commit(dbname, collname string)

	SetData(title string, starscount int)
}