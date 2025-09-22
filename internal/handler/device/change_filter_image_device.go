package device

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/device"
	"xinde/internal/handler/common"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// ChangeFilterImageDevice handles changing the device type of a filter image config.
// @Summary      更改所属设备
// @Description  将一条筛选条件图片配置移动到另一个设备类型下
// @Tags         FilterImage
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "筛选条件图片配置 ID"
// @Param        body body      dto.ChangeDeviceTypeReq true "请求体"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "操作成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      404 {object} response.Response "配置或目标设备不存在"
// @Failure      409 {object} response.Response "目标设备下已存在相同的筛选值配置 (Conflict)"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/filter_image/change/device_type/{id} [patch]
func (ctrl *Controller) ChangeFilterImageDevice(c *gin.Context) {

	id, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorFilterImageIDInvalid)
		logger.Error("/admin/filter_image/change/device_type 无效的筛选下拉列表图片ID格式: " + err.Error())
		return
	}

	var req *dto.ChangeDeviceTypeReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "参数绑定错误: "+err.Error())
		logger.Error("/admin/filter_image/change/device_type 参数绑定错误: " + err.Error())
		return
	}

	// 剩余的工作交由service层处理
	err = ctrl.service.ChangeFilterImageDevice(id, req.DeviceTypeID)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorFilterImageValueConflict:
			response.Error(c, http.StatusConflict, response.CodeConflict, stderr.ErrorFilterImageValueConflict)
		case stderr.ErrorFilterImageNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorFilterImageNotFound)
		case stderr.ErrorDeviceNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorDeviceNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "更新下拉列表图片设备发生错误: "+err.Error())
			logger.Error("/admin/filter_image/change/device_type 更新下拉列表图片设备发生错误: " + err.Error())
		}
		return
	}
	response.Success(c, nil)
}
