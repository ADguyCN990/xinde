package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xinde/pkg/stderr"
)

// getIDFromUrl 从类似/admin/account/:id这种格式中提取userID
func (ctrl *Controller) getIDFromUrl(c *gin.Context) (uint, error) {
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
