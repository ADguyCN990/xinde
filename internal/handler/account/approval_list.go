package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	dto "xinde/internal/dto/account"
	model "xinde/internal/model/account"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// ApprovalList handles approval user list.
// @Summary 管理员查看用户审批列表
// @Description 根据前端传来的字段，返回对应的待审批列表/已同意申请/已拒绝申请的用户列表，根据创建时间降序排列
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dto.ApprovalListReq true "审批列表 Request"
// @Success 200 {object} dto.ApprovalListResp "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/account/approval/list [get]
func (ctrl *Controller) ApprovalList(c *gin.Context) {
	var req dto.ApprovalListReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		logger.Error("admin/account/approval/list 参数绑定错误: " + err.Error())
		return
	}

	// 优先取前端传过来的pageSize，如果没有取默认值
	if req.PageSize == 0 {
		req.PageSize = viper.GetInt("page.defaultPageSize")
	}

	// 对前端传过来的status做校验
	var status int
	switch req.Status {
	case "pending":
		status = model.UserPending
	case "":
		status = model.UserPending
	case "approved":
		status = model.UserApproved
	case "rejected":
		status = model.UserRejected
	default:
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "输入正确的status值，status只包含 pending/approved/rejected 这三种情况")
		logger.Error("admin/account/approval/list 非法的status值: " + req.Status)
		return
	}

	// 参数校验完毕，剩余的工作交由Service层处理
	list, err := ctrl.accountService.GetApprovalUserList(req.Page, req.PageSize, status)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorDbNil:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
			logger.Error("admin/account/approval/list " + err.Error())
		case stderr.ErrorOverLargePage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至最后一页", stderr.ErrorOverLargePage), list)
		case stderr.ErrorOverSmallPage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至第一页", stderr.ErrorOverSmallPage), list)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
			logger.Error("admin/account/approval/list " + err.Error())
		}
		return
	}

	response.Success(c, list)
}
