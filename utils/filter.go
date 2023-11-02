package utils

import (
	"regexp"

	"github.com/labstack/gommon/log"
	"github.com/samber/lo"
)

func FilterByRegexp(regexpList []string, target string) bool {
	return lo.ContainsBy(regexpList, func(filter string) bool {
		re, err := regexp.Compile(filter)
		if err != nil {
			log.Warnf("error filter regexp %v", filter)
			return false
		}
		return re.Match([]byte(target))
	})
}
