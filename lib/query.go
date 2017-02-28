package lib

import (
	"fmt"
	"strconv"
)

// Query is all all query options
type Query struct {
	Q   string
	Opt *Option
}

func (q *Query) String() string {
	q.build()
	return fmt.Sprintf("%s", q.Q)
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
