package dockerregistry

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/tcaty/update-watcher/internal/core"
	"github.com/tcaty/update-watcher/pkg/markdown"
)

func splitNsAndRepo(image string) (string, string, error) {
	b := []byte(image)
	i := bytes.IndexByte(b, byte('/'))

	if i < 0 {
		return "", "", errors.New("docker image should fit the format {namespace}/{repository}")
	}

	ns, repo := string(b[:i]), string(b[i+1:])
	return ns, repo, nil
}

func createHrefs(vrs core.VersionRecords) []fmt.Stringer {
	hrefs := make([]fmt.Stringer, 0)
	for t, v := range vrs {
		text := fmt.Sprintf("%s:%s", t, v)
		link := fmt.Sprintf("https://hub.docker.com/r/%s/tags", t)
		href := markdown.NewHref(text, link)
		hrefs = append(hrefs, href)
	}
	return hrefs
}
