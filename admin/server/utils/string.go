// @program:     flashbear
// @file:        string.go.go
// @author:      ac
// @create:      2024-10-31 10:51
// @description:
package utils

import (
	"regexp"
)

func GetShareId(site string) string {
	re := regexp.MustCompile(`share_id=([a-zA-Z0-9]+)`)
	match := re.FindStringSubmatch(site)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
