package device

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	dto "xinde/internal/dto/device"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// FilterImageList handles fetching a paginated list of filter images.
// @Summary      获取筛选条件图片列表
// @Description  分页获取已配置的筛选条件图片列表，可按设备类型筛选
// @Tags         FilterImage
// @Accept       json
// @Produce      json
// @Param        page query int false "页码"
// @Param        page_size query int false "每页数量"
// @Param        device_type_id query int true "设备类型ID (用于筛选)"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response{data=dto.ListPageData} "成功返回列表"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/filter_image/list [get]
func (ctrl *Controller) FilterImageList(c *gin.Context) {
	var req *dto.FilterImageListReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
		logger.Error("/admin/filter_image/list 绑定参数错误: " + err.Error())
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = viper.GetInt("page.defaultPageSize")
	}

	// 将剩余的工作交由service处理
	list, err := ctrl.service.FilterImageList(req.DeviceTypeID, req.Page, req.PageSize)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorOverLargePage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至最后一页", stderr.ErrorOverLargePage), list)
		case stderr.ErrorOverSmallPage:
			response.SuccessWithMessage(c, fmt.Sprintf("%s, 跳转至第一页", stderr.ErrorOverSmallPage), list)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "获取列表失败")
			logger.Error("/admin/filter_image/list 获取列表失败: " + err.Error())
		}
		return
	}
	response.Success(c, list)
}
