package group

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	dto "xinde/internal/dto/group"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// List handles fetching a paginated flat list of groups for the backend.
// @Summary      获取后台分组列表
// @Description  获取一个扁平化、带分页的分组列表，用于后台管理
// @Tags         Group
// @Accept       json
// @Produce      json
// @Param        page query int false "页码 (默认: 1)"
// @Param        page_size query int false "每页数量 (默认: from config)"
// @Security     ApiKeyAuth
// @Success      200 {object} dto.ListResp "成功返回分组列表"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/group/list [get]
func (ctrl *Controller) List(c *gin.Context) {
	var req *dto.ListReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
		logger.Error("/admin/group/list 绑定参数错误: " + err.Error())
		return
	}

	if req.PageSize == 0 {
		req.PageSize = viper.GetInt("page.defaultPageSize")
	}

	// 参数校验完毕，剩余的工作交由Service处理
	list, err := ctrl.Service.GetGroupList(req.Page, req.PageSize)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorDbNil:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/group/list " + err.Error())
		case stderr.ErrorOverLargePage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至最后一页", stderr.ErrorOverLargePage), list)
		case stderr.ErrorOverSmallPage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至第一页", stderr.ErrorOverSmallPage), list)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/group/list 获取分组列表失败 " + err.Error())
		}
		return
	}
	response.Success(c, list)
}
