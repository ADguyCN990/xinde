package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/account"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// UpdatePassword handles admin update user's password.
// @Summary 修改密码
// @Description 管理员根据用户ID修改用户的密码
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dto.UpdatePasswordReq true "UpdatePassword　Request"
// @Success 200 {object} response.Response "修改密码成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "access_token有错误"
// @Failure 403 {object} response.Response "没有管理员权限"
// @Failure 404 {object} response.Response "没有该用户"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/account/password/{id} [patch]
func (ctrl *Controller) UpdatePassword(c *gin.Context) {
	id, err := ctrl.getIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorUserIDInvalid)
		logger.Error("/admin/account/password/patch 无效的用户ID格式: " + c.Param("id"))
		return
	}

	var req dto.UpdatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		logger.Error("/admin/account/password/patch 参数绑定错误: " + err.Error())
		return
	}

	// 完成参数校验，将剩余的工作交由service处理
	err = ctrl.accountService.UpdatePassword(id, req.Password)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorUserNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorUserNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error(fmt.Sprintf("/admin/account/password/patch 修改用户密码失败! 用户ID: %d 错误: %s", id, err.Error()))
		}
		return
	}

	response.Success(c, nil)
}
