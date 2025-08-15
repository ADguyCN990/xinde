package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	dto "xinde/internal/dto/account"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// List handles user list.
// @Summary 管理员查看用户列表
// @Description 返回已被通过注册申请的用户信息
// @Tags Account
// @Accept json
// @Produce json
// @Param page query int false "当前页数，可选，默认为1"
// @Param page_size query int false "一页的内容数量，可选，默认为设置的默认值"
// @Success 200 {object} dto.ListResp "查询成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/admin/account/list [get]
func (ctrl *Controller) List(c *gin.Context) {
	var req dto.ListReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeUnauthorized, err.Error())
		logger.Error("account/list 绑定参数错误: " + err.Error())
		return
	}

	if req.PageSize == 0 {
		req.PageSize = viper.GetInt("page.defaultPageSize")
	}

	// 参数校验完毕，剩余的工作交由Service层处理
	list, err := ctrl.accountService.GetUserList(req.Page, req.PageSize)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorDbNil:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
			logger.Error("admin/account/list " + err.Error())
		case stderr.ErrorOverLargePage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至最后一页", stderr.ErrorOverLargePage), list)
		case stderr.ErrorOverSmallPage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至第一页", stderr.ErrorOverSmallPage), list)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
			logger.Error("admin/account/list " + err.Error())
		}
		return
	}
	response.Success(c, list)
}
