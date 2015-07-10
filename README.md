# githubstars

Work in progress

## Install
```go get github.com/saromanov/githubstars```

## Usage
```go
package main
import
(
	"github.com/saromanov/githubstars"
)

func main() {
	res := githubstars.Init()
	res.Show(githubstars.Options{Numstars: ">2000", Language: "go"})
}
```

## API

### githubstars.Init()

### Show(opt Options)
Output results

## Commit(name string)
Store results in mongo with namse as collection.
If name is empty, results will store in collection with default name

### CompareWith(dbtitle string)
Compare results with one of the prevous results

### AvailableResults(opt Options)
Get all of names with this query

