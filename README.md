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

