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

func (h *Href) Get() string {
	return fmt.Sprintf("[%s](%s)", h.text, h.link)
}
