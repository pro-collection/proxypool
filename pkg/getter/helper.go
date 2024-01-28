package getter

import (
	"regexp"
	"strconv"
)

func extractSpeed(oritext string) int64 {
	reg := regexp.MustCompile(`\[1-9\]\d\*\\.\?\d\*`)
	temp := reg.FindString(oritext)
	if temp != "" {
		speed, _ := strconv.ParseInt(temp, 10, 64)
		return speed
	}
	return -1
}
