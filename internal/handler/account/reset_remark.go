package account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

func (ctrl *Controller) ResetRemark(c *gin.Context) {
	id, err := ctrl.getIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorUserIDInvalid)
		logger.Error("/admin/account/reset/remark/ 无效的用户ID格式: " + c.Param("id"))
		return
	}

	// 无需参数校验，将剩余的工作交由service处理
	ctrl.accountService.ResetRemark(id)
}
