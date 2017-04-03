package lib

import (
	"fmt"
	"strconv"
	"strings"
)

// Query is all all query options
type Query struct {
	Q   string
	Opt *Option

	// Total of params contained in the query
	NbParams uint32
	Values   []interface{}
}

// NewQuery initiates a query
func NewQuery(query string, opt *Option, nbParams uint32) *Query {
	return &Query{
		Q:        query,
		Opt:      opt,
		NbParams: nbParams,
		Values:   make([]interface{}, 0),
	}
}

func (q *Query) String() string {
	q.build()
	return fmt.Sprintf("%s", q.Q)
}

// SetFilters adds sql filters (LIKE) to the query
func (q *Query) SetFilters(filters ...string) {
	for i := 0; i < len(filters)-1; i += 2 {
		queryFilter := filters[i]
		filter := q.Opt.GetFilter(queryFilter)
		if filter != nil {
			switch queryFilter {
			case SEARCH:
				searchFilters := strings.Split(filters[i+1], "|")
				q.Q += " AND ("
				for _, searchFilter := range searchFilters {
					q.NbParams++
					q.Q += "LOWER(" + searchFilter + ") LIKE $" + fmt.Sprintf("%d", q.NbParams) + " OR "
					q.Values = append(q.Values, "%"+filter.Value+"%")
				}
				q.Q = q.Q[0 : len(q.Q)-4] // Remove last 'OR'
				q.Q += ")"
				break
			default:
				q.NbParams++
				q.Q += " AND " + filters[i+1] + " = $" + fmt.Sprintf("%d", q.NbParams)
				q.Values = append(q.Values, filter.Value)
			}
		}
	}
}

// AddValues add the values for the SQL query
func (q *Query) AddValues(values ...interface{}) {
	q.Values = append(q.Values, values...)
}

func (q *Query) build() {
	if q.Opt == nil {
		return
	}
	var str string
	if q.Opt.Limit > 0 {
		str = strconv.FormatInt(q.Opt.Limit, 10)
		q.Q += " LIMIT " + str
	}
	if q.Opt.Offset > 0 {
		str = strconv.FormatInt(q.Opt.Offset, 10)
		q.Q += " OFFSET " + str
	}
	q.Q += ";"
}
