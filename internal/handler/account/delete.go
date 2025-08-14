package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "xinde/internal/dto/account"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// DeleteUser handles admin approve.
// @Summary 删除用户
// @Description 根据ID删除用户
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "access_token有错误"
// @Failure 403 {object} response.Response "没有管理员权限"
// @Failure 404 {object} response.Response "没有该用户"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/account/{id} [delete]
func (ctrl *Controller) DeleteUser(c *gin.Context) {
	id, err := ctrl.getIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorUserIDInvalid)
		logger.Error("/admin/account/delete 无效的用户ID格式: " + c.Param("id"))
		return
	}

	// 无需参数校验，将剩余的工作交由service处理
	err = ctrl.accountService.DeleteUser(id)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorUserNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorUserNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error(fmt.Sprintf("/admin/account/delete/ 删除用户失败! 用户ID: %d 错误: %s", id, err.Error()))
		}
		return
	}

	response.Success(c, nil)
}
