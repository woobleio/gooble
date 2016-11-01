package util

import (
  "fmt"
  "strconv"
  "strings"

  "github.com/gin-gonic/gin"
)

type Query struct {
  Q    string
  Opt  *Option
}

func (q *Query) String() string {
  q.build()
  return fmt.Sprintf("%s", q.Q)
}

func (q *Query) build() {
  var str string
  if q.Opt.limit > 0 {
    str = strconv.FormatInt(q.Opt.limit, 10)
    q.Q += " LIMIT " + str
  }
  if q.Opt.offset > 0 {
    str = strconv.FormatInt(q.Opt.offset, 10)
    q.Q += " OFFSET " + str
  }
  q.Q += ";"
}

type Option struct {
  limit     int64
  offset    int64
  Populate  []string
}

func ParseOptions(c *gin.Context) (Option) {
  pLimit := c.DefaultQuery("limit", "0")
  pOffset := c.DefaultQuery("offset", "0")
  pPopulate := c.Query("populate")

  limit, errLimit := strconv.ParseInt(pLimit, 10, 64)
  if errLimit != nil {
    limit = 0
    fmt.Errorf("Option limit(%s) is invalid, error : %s", pLimit, errLimit)
  }

  offset, errOffset := strconv.ParseInt(pOffset, 10, 64)
  if errOffset != nil {
    offset = 0
    fmt.Errorf("Option offset(%s) is invalid, error : %s", pOffset, errOffset)
  }

  // TODO is populate needed
  populate := strings.Split(pPopulate, ",")

  return Option{limit, offset, populate}
}
