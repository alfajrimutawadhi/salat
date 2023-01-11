package common

import (
	"fmt"
	"strconv"
	"strings"
)

// t format HH:MM
func ConvTime24To12(t string) string {
	adverb := "AM"
	tmp := strings.Split(t, ":")
	h, err := strconv.Atoi(tmp[0])
	if err != nil {
		HandleError(err)
	}
	if h > 12 {
		h -= 12
		tmp[0] = fmt.Sprintf("0%d", h)
		adverb = "PM"
	} else if h == 12 {
		adverb = "PM"
	}
	return fmt.Sprintf("%s:%s%s", tmp[0], tmp[1], adverb)
}
