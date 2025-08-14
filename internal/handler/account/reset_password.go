package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// ResetPassword handles admin reset user's password.
// @Summary 重置用户密码
// @Description 管理员根据用户ID重置用户的密码为123456
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "重置密码成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "access_token有错误"
// @Failure 403 {object} response.Response "没有管理员权限"
// @Failure 404 {object} response.Response "没有该用户"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/account/reset/password/{id} [post]
func (ctrl *Controller) ResetPassword(c *gin.Context) {
	id, err := ctrl.getIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorUserIDInvalid)
		logger.Error("/admin/account/reset/password/ 无效的用户ID格式: " + c.Param("id"))
		return
	}

	// 无需参数校验，将剩余的工作交给Service处理
	err = ctrl.accountService.ResetPassword(id)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorUserNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorUserNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error(fmt.Sprintf("/admin/account/reset/password/ 重置用户密码失败! 用户ID: %d 错误: %s", id, err.Error()))
		}
		return
	}

	response.Success(c, nil)
}
