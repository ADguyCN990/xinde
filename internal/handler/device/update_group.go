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

// UpdateGroup handles changing the group of a device type.
// @Summary      更换设备类型的分组
// @Description  将一个设备类型移动到另一个指定的分组下
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "设备类型 ID"
// @Param        body body     dto.ChangeGroupReq true "包含新分组ID的请求体"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "更换成功"
// @Failure      400 {object} response.Response "请求参数错误或无效ID"
// @Failure      404 {object} response.Response "设备类型或目标分组不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/update/group/{id} [patch]
func (ctrl *Controller) UpdateGroup(c *gin.Context) {
	deviceTypeID, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorDeviceIDInvalid)
		logger.Error("/admin/device/update/group 无效的设备类型ID格式: " + err.Error())
		return
	}

	var req *dto.ChangeGroupReq
	if err = c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
		logger.Error("/admin/device/update/group 绑定参数错误: " + err.Error())
		return
	}

	// 剩余的工作交由service处理
	err = ctrl.service.UpdateGroup(deviceTypeID, req.GroupID)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorGroupNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorGroupNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/device/update/group 更新设备类型分组错误: " + err.Error())
		}
		return
	}
	response.Success(c, nil)
}
