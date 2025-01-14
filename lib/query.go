package lib

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Query are query options
type Query struct {
	Q   string
	Opt *Option

	// Total of params contained in the query
	Values []interface{}
}

// NewQuery initiates a query
func NewQuery(query string, opt *Option) *Query {
	return &Query{
		Q:      query,
		Opt:    opt,
		Values: make([]interface{}, 0),
	}
}

func (q *Query) String() string {
	q.build()
	return fmt.Sprintf("%s", q.Q)
}

// SetFilters adds sql filters (LIKE) to the query
// isAnding tells the parser if it has to start whith a WHERE clause (false) or AND close (true)
// filters are the string values to filter with LIKE %filtre%
func (q *Query) SetFilters(isAnding bool, filters ...string) {
	for i := 0; i < len(filters)-1; i += 2 {
		queryFilter := filters[i]
		filter := q.Opt.GetFilter(queryFilter)
		if filter != nil {
			if i == 0 && !isAnding {
				q.Q += " WHERE "
			} else {
				q.Q += " AND "
			}
			switch queryFilter {
			case SEARCH:
				searchFilters := strings.Split(filters[i+1], "|")
				q.Q += "("
				for _, searchFilter := range searchFilters {
					q.Values = append(q.Values, "%"+filter.Value+"%")
					q.Q += "LOWER(" + searchFilter + ") LIKE LOWER($" + fmt.Sprintf("%d", len(q.Values)) + ") OR "
				}
				q.Q = q.Q[0 : len(q.Q)-4] // Remove last 'OR'
				q.Q += ")"
				break
			default:
				q.Values = append(q.Values, filter.Value)
				q.Q += filters[i+1] + " = $" + fmt.Sprintf("%d", len(q.Values))
			}
		}
	}
}

// SetBulkInsert adds bulk insert to the query
// baseValues are values that doesn't change or bulk inserting (can be empty)
// atrs are attributes of the values interface
// values are the dynamic values to bulk insert (attrs relates to this interfaces slice)
func (q *Query) SetBulkInsert(baseValues []interface{}, attrs []string, values ...interface{}) {
	var index = 1
	for _, v := range values {
		q.Q += ` (`
		for _, base := range baseValues {
			q.Q += `$` + fmt.Sprintf("%d", index) + `,`
			q.Values = append(q.Values, base)
			index++
		}

		for _, attr := range attrs {
			q.Q += `$` + fmt.Sprintf("%d", index) + `,`
			q.Values = append(q.Values, reflect.ValueOf(v).FieldByName(attr).Interface())
			index++
		}
		q.Q = strings.TrimRight(q.Q, ",") + `),`
	}
	q.Q = strings.TrimRight(q.Q, ",")
}

// SetOrder adds order to the query
func (q *Query) SetOrder(key string, field string) {
	sort := q.Opt.GetSort(key)
	if sort != nil {
		q.Q += " ORDER BY " + field + " " + sort.Order
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
