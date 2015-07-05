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
	res.Get(">2000", "","go")
}
```

