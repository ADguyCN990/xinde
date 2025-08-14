package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	dto "xinde/internal/dto/account"
	model "xinde/internal/model/account"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// Approve handles admin approve.
// @Summary 批准用户申请
// @Description 批准用户申请，管理员决定是否同意用户的注册申请
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dto.ApproveReq true "Login Request"
// @Success 200 {object} dto.ApproveResp "审批成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "access_token有错误"
// @Failure 403 {object} response.Response "没有管理员权限"
// @Failure 404 {object} response.Response "没有该用户"
// @Failure 409 {object} response.Response "用户已经被审批"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/account/approval/{id} [post]
func (ctrl *Controller) Approve(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorUserIDInvalid)
		logger.Error("admin/account/approval 无效的用户ID格式: " + idStr)
		return
	}
	if id < 1 {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorUserIDInvalid)
		logger.Error("admin/account/approval 无效的用户ID格式: " + idStr)
		return
	}

	var req *dto.ApproveReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		logger.Error("admin/account/approval 参数绑定错误: " + err.Error())
		return
	}

	// 将前端传来的status参数转换成对应的`is_user`字段
	statusMap := map[string]int{
		"approve": model.UserApproved,
		"reject":  model.UserRejected,
	}
	status, exists := statusMap[req.Status]
	if !exists {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "status只能是approve或reject")
		logger.Error("admin/account/approval 非法status: " + req.Status)
		return
	}

	// 参数校验完毕，剩余的工作交由service处理
	err = ctrl.accountService.ApproveUser(uint(id), status, req.Why)

	// 根据错误，向前端返回不同的响应
	if err != nil {
		ctrl.handleApproveError(c, err, uint(id))
		return
	}

	response.Success(c, nil)
}

// handleApproveError 错误处理逻辑提取为单独方法
func (ctrl *Controller) handleApproveError(c *gin.Context, err error, userID uint) {
	errorMap := map[string]struct {
		status int
		code   int
		msg    string
	}{
		stderr.ErrorUserBanned:   {http.StatusConflict, response.CodeConflict, "用户已被拒绝，请勿重复审批"},
		stderr.ErrorUserPassed:   {http.StatusConflict, response.CodeConflict, "用户已被通过，请勿重复审批"},
		stderr.ErrorUserNotFound: {http.StatusNotFound, response.CodeNotFound, "用户不存在"},
	}

	if errInfo, exists := errorMap[err.Error()]; exists {
		response.Error(c, errInfo.status, errInfo.code, errInfo.msg)
	} else {
		logger.Error(fmt.Sprintf("admin/account/approval 审批失败! 用户ID: %d, 错误: %s", userID, err.Error()))
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "服务器内部错误")
	}
}
