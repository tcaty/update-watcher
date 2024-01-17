package markdown

import "fmt"

func CreateUnorderedList(arr []string) string {
	res := ""
	for _, v := range arr {
		res = fmt.Sprintf("%s* %s\n", res, v)
	}
	return res
}
