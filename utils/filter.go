package utils

import (
	"log"
	"regexp"

	"github.com/samber/lo"
)

func FilterByRegexp(regexpList []string, target string) bool {
	return lo.ContainsBy(regexpList, func(filter string) bool {
		re, err := regexp.Compile(filter)
		if err != nil {
			log.Printf("error filter regexp %v\n", filter)
			return false
		}
		return re.Match([]byte(target))
	})
}
