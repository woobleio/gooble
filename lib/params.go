package lib

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Option is query option
type Option struct {
	Search string
	Limit  int64
	Offset int64
}

// ParseOptions parses query options
func ParseOptions(c *gin.Context) Option {
	pPage := c.DefaultQuery("page", "1")
	pPerPage := c.DefaultQuery("perPage", "15")
	pSearch := c.DefaultQuery("search", "")

	page, errPage := strconv.ParseInt(pPage, 10, 64)
	if errPage != nil || page <= 0 {
		page = 1
	}

	perPage, errPerPage := strconv.ParseInt(pPerPage, 10, 64)
	if errPerPage != nil || perPage <= 0 {
		perPage = 15
	}

	return Option{pSearch, perPage, ((page - 1) * perPage)}
}
