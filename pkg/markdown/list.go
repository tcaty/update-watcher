package markdown

import "fmt"

func CreateUnorderedList(slice []fmt.Stringer) string {
	res := ""
	for _, v := range slice {
		res = fmt.Sprintf("%s* %s\n", res, v.String())
	}
	return res
}
