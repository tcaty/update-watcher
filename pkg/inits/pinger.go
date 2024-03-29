package inits

import (
	"github.com/tcaty/update-watcher/pkg/utils"
)

type Pinger interface {
	Ping() error
}

// like PingAll but convert to Pinger type before pinging
// and convert to initial slice type after that
func PingGeneric[T comparable](s []T) error {
	ps := utils.ConvertSlice[T, Pinger](s)
	return Ping(ps)
}

func Ping(ps []Pinger) error {
	for _, p := range ps {
		if err := p.Ping(); err != nil {
			return err
		}
	}
	return nil
}
