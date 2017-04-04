package lib

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Option is query option
type Option struct {
	Filters *[]Filter
	Limit   int64
	Offset  int64
	Sort    *Sort
}

// Sort defines a ordering
type Sort struct {
	Field string
	Order string
}

// Filter is a query filter
type Filter struct {
	ID    string
	Value string
}

// Filters
const (
	CREATOR = "creator"
	SEARCH  = "search"
)

var filters = []string{
	CREATOR,
	SEARCH,
}

// ParseOptions parses query options
func ParseOptions(c *gin.Context) Option {
	filtersObj := make([]Filter, 0)
	sort := c.DefaultQuery("sort", "")
	pPage := c.DefaultQuery("page", "1")
	pPerPage := c.DefaultQuery("perPage", "15")
	for _, filter := range filters {
		qFilter := c.DefaultQuery(filter, "")
		if qFilter != "" {
			filtersObj = append(filtersObj, Filter{filter, qFilter})
		}
	}

	var sortObj *Sort
	if sort != "" {
		switch sort[0] {
		case '-':
			sortObj = &Sort{sort[1:len(sort)], "DESC"}
		default:
			sortObj = &Sort{sort, "ASC"}
		}
	}

	page, errPage := strconv.ParseInt(pPage, 10, 64)
	if errPage != nil || page <= 0 {
		page = 1
	}

	perPage, errPerPage := strconv.ParseInt(pPerPage, 10, 64)
	if errPerPage != nil || perPage <= 0 {
		perPage = 15
	}

	return Option{&filtersObj, perPage, ((page - 1) * perPage), sortObj}
}

// GetFilter return true if the options contains the filter
func (o *Option) GetFilter(queryFilter string) *Filter {
	for _, filter := range *o.Filters {
		if filter.ID == queryFilter {
			return &filter
		}
	}
	return nil
}
