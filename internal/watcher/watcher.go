package watcher

import "fmt"

type Watcher interface {
	GetLastVersion() (string, error)
}

func Tick(w Watcher) {
	fmt.Println(w.GetLastVersion())
}
