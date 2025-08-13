package account

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	dto "xinde/internal/dto/account"
	"xinde/pkg/logger"
	"xinde/pkg/response"
)

// Approval List handles approval user list.
// @Summary 管理员查看用户审批列表
// @Description 根据前端传来的字段，返回对应的待审批列表/已同意申请/已拒绝申请的用户列表，根据创建时间降序排列
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dto.ApprovalListReq true "审批列表 Request"
// @Success 200 {object} dto.ListResp "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/account/approval/list [get]
func (ctrl *Controller) approvalList(c *gin.Context) {
	var req dto.ApprovalListReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error())
		logger.Error("account/approval/list 参数绑定错误: " + err.Error())
		return
	}

	if req.PageSize == 0 {
		req.PageSize = viper.GetInt("page.defaultPageSize")
	}

	//
}
