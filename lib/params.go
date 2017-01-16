package lib

import (
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
)

// Option is query option
type Option struct {
	Limit  int64
	Offset int64
}

// ParseOptions parses query options
func ParseOptions(c *gin.Context) Option {
	pPage := c.DefaultQuery("page", "1")
	pPerPage := c.DefaultQuery("perPage", "15")

	page, errPage := strconv.ParseInt(pPage, 10, 64)
	if errPage != nil || page <= 0 {
		page = 1
	}

	perPage, errPerPage := strconv.ParseInt(pPerPage, 10, 64)
	if errPerPage != nil || perPage <= 0 {
		perPage = 15
	}

	return Option{perPage, ((page - 1) * perPage)}
}
