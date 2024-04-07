package response

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	HeaderTotalCount    = "X-Total-Count"
	HeaderCurrentOffset = "X-Offset"
	HeaderLimit         = "X-Limit"
	HeaderNextPage      = "X-Next-Page"
)

func AddPaginationHeaders(c *gin.Context, offset, limit int, totalCnt uint64, nextPage bool) {
	c.Header(HeaderTotalCount, fmt.Sprintf("%d", totalCnt))
	c.Header(HeaderCurrentOffset, fmt.Sprintf("%d", offset))
	c.Header(HeaderLimit, fmt.Sprintf("%d", limit))
	c.Header(HeaderNextPage, fmt.Sprintf("%v", nextPage))
}
