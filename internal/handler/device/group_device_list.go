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

// GroupDeviceList handles fetching a list of device types.
// @Summary      根据分组获取设备类型列表
// @Description  分页获取设备类型列表，用于前台展示
// @Tags         group
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "设备类型 ID"
// @Security     ApiKeyAuth
// @Success      200 {object} dto.GroupDeviceListResp "成功返回列表"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      404 {object} response.Response "分组不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/group/device/list [get]
func (ctrl *Controller) GroupDeviceList(c *gin.Context) {
	groupID, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorGroupIDInvalid)
		logger.Error("/admin/device/group/list 无效的分组ID格式: " + err.Error())
		return
	}
	_ = dto.ListPageData{}
	// 剩余的工作交由service层处理
	list, err := ctrl.service.GroupDeviceList(groupID)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorGroupNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorGroupNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/device/group/list 根据分组ID查找设备出错: " + err.Error())
		}
		return
	}
	response.Success(c, list)
}
