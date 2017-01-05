package lib

import (
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
)

// Option is query option
type Option struct {
	limit  int64
	offset int64
}

// ParseOptions parses query options
func ParseOptions(c *gin.Context) Option {
	pLimit := c.DefaultQuery("numberOf", "0")
	pOffset := c.DefaultQuery("to", "0")

	limit, errLimit := strconv.ParseInt(pLimit, 10, 64)
	if errLimit != nil {
		limit = 0
		// fmt.Printf("Option limit(%s) is invalid, error : %s", pLimit, errLimit)
	}

	offset, errOffset := strconv.ParseInt(pOffset, 10, 64)
	if errOffset != nil {
		offset = 0
		// fmt.Printf("Option offset(%s) is invalid, error : %s", pOffset, errOffset)
	}

	return Option{limit, offset}
}
