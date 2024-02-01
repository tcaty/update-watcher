package inits

import (
	"log/slog"

	"github.com/tcaty/update-watcher/pkg/utils"
)

type Toggler interface {
	Enabled() bool
}

// like ExcludeDisabled but convert to Toggler type before excluding
// and convert to initial slice type after that
func ExcludeDisabledGeneric[T comparable](s []T, log bool) []T {
	tgs := utils.ConvertSlice[T, Toggler](s)
	filtered := ExcludeDisabled(tgs, log)
	return utils.ConvertSlice[Toggler, T](filtered)
}

func ExcludeDisabled(tgs []Toggler, log bool) []Toggler {
	res := make([]Toggler, 0)
	for _, tg := range tgs {
		if tg.Enabled() {
			res = append(res, tg)
		}
		if log {
			slog.Debug(
				"initializing",
				"toggler", tg,
				"enabled", tg.Enabled(),
			)
		}
	}
	return res
}
