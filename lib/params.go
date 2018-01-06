package lib

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Option is query option
type Option struct {
	Filters   *[]Filter
	Limit     int64
	Offset    int64
	Sort      *Sort
	Populates []string
}

// Sort defines a ordering
type Sort struct {
	Key   string
	Order string
}

// Filter is a query filter
type Filter struct {
	ID    string
	Value string
}

// Filters and orders
const (
	CREATED_AT = "createdAt"
	CREATOR    = "creator"
	SEARCH     = "search"
	NB_USE     = "nbUse"
)

var filters = []string{
	CREATED_AT,
	CREATOR,
	SEARCH,
	NB_USE,
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
		case '+':
			sort = sort[1:len(sort)]
			fallthrough
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

	populates := strings.Split(c.DefaultQuery("populate", ""), ",")

	return Option{&filtersObj, perPage, ((page - 1) * perPage), sortObj, populates}
}

// GetFilter returns the filter object
func (o *Option) GetFilter(queryFilter string) *Filter {
	for _, filter := range *o.Filters {
		if filter.ID == queryFilter {
			return &filter
		}
	}
	return nil
}

// GetSort returns the sort object
// It controls if the sort key exists in the filters, for security reasons
// so that it only consider known fields
func (o *Option) GetSort(orderKey string) *Sort {
	if o.Sort == nil {
		return nil
	}
	for _, filter := range filters {
		if filter == orderKey && o.Sort.Key == orderKey {
			return o.Sort
		}
	}
	return nil
}

// HasPopulate return true if it contains populate "key"
func (o *Option) HasPopulate(key string) bool {
	for _, populate := range o.Populates {
		if populate == key {
			return true
		}
	}
	return false
}
