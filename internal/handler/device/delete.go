package device

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xinde/internal/handler/common"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// Delete handles the complete deletion of a device type and its related data.
// @Summary      删除设备类型
// @Description  彻底删除一个设备类型及其所有方案、主图和导入的Excel文件
// @Tags         Device
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "设备类型 ID"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "删除成功"
// @Failure      400 {object} response.Response "请求参数错误或无效ID"
// @Failure      404 {object} response.Response "设备类型不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/device/delete/{id} [delete]
func (ctrl *Controller) Delete(c *gin.Context) {
	deviceTypeID, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeForbidden, stderr.ErrorDeviceIDInvalid)
		logger.Error("/admin/device/delete 无效的设备类型ID格式: " + err.Error())
		return
	}

	// 剩余的工作交由service处理
	err = ctrl.service.Delete(deviceTypeID)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorDeviceNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorDeviceNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/device/delete 删除设备类型失败: " + err.Error())
		}
		return
	}
	response.Success(c, nil)
}
