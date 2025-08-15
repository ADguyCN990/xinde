package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xinde/pkg/stderr"
)

func GetIDFromUrl(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	if id < 1 {
		return 0, fmt.Errorf(stderr.ErrorUserIDInvalid)
	}
	return uint(id), nil
}
