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

// UpdateName handles changing the name of a device type.
// @Summary      更换设备类型的名称
// @Description  更新一个设备类型的名称
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "设备类型 ID"
// @Param        body body      dto.ChangeNameReq true "包含新名称的请求体"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "更换成功"
// @Failure      400 {object} response.Response "请求参数错误或无效ID"
// @Failure      404 {object} response.Response "设备类型不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/device/update/name/{id} [patch]
func (ctrl *Controller) UpdateName(c *gin.Context) {
	deviceTypeID, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorDeviceIDInvalid)
		logger.Error("/admin/device/update/name 无效的设备类型ID格式: " + err.Error())
		return
	}

	var req *dto.ChangeNameReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
		logger.Error("/admin/device/update/name 绑定参数错误: " + err.Error())
		return
	}

	// 将剩余的工作交由service处理
	err = ctrl.service.UpdateName(deviceTypeID, req.Name)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "更改设备名称类型名称失败: "+err.Error())
		logger.Error("/admin/device/update/name 更改设备名称类型失败: " + err.Error())
		return
	}
	response.Success(c, nil)
}
