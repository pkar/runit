# log

That's it, just simple leveled logging with the Go std lib log.


```go
import (
	"github.com/pkar/log"
)

func main() {
	// 0 debug
	// 1 info
	// 2 warn
	// 3 error
	// 4 fatal
	log.SetLevel(0)

	//log.SetOutput(someiowriter) // default os.Stdout

	log.Debug("debug")
	log.Debugf("debug %s", "d")
	log.Info("info")
	log.Infof("info %s", "i")
	log.Print("line")
	log.Println("line")
	log.Printf("print %s", "p")
	log.Warn("warn")
	log.Warnf("warn %v", "w")
	log.Error("error")
	log.Errorf("error %v", "e")
	log.Fatal("bye")
	log.Fatalf("bye %s", "bye")
}
```
