package markdown

import "fmt"

type Href struct {
	text string
	link string
}

func NewHref(text string, link string) *Href {
	return &Href{
		text: text,
		link: link,
	}
}

// Sprint creates the markdown href and returns the resulting string.
func (h *Href) Sprint() string {
	return fmt.Sprintf("[%s](%s)", h.text, h.link)
}
