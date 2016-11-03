package lib

import (
  "fmt"
  "strconv"

  "github.com/gin-gonic/gin"
)

type Option struct {
  limit     int64
  offset    int64
}

func ParseOptions(c *gin.Context) (Option) {
  pLimit := c.DefaultQuery("limit", "0")
  pOffset := c.DefaultQuery("offset", "0")

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

  return Option{limit, offset}
}
