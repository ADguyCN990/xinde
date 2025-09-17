package device

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xinde/internal/handler/common"
	"xinde/internal/middleware/auth"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// UpdateImport handles re-importing devices from an Excel file for an existing DeviceType.
// @Summary      更新导入设备方案
// @Description  为一个已存在的设备类型上传新的Excel文件，覆盖其下所有方案
// @Tags         Device
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   path      int  true  "设备类型 ID"
// @Param        device formData file true "包含新设备方案的Excel文件"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "更新导入成功"
// @Failure      400 {object} response.Response "请求参数错误或无效ID"
// @Failure      404 {object} response.Response "设备类型不存在"
// @Failure      500 {object} response.Response "服务器内部错误或导入失败"
// @Router       /api/v1/admin/device/import/{id} [put]
func (ctrl *Controller) UpdateImport(c *gin.Context) {
	id, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorDeviceIDInvalid)
		logger.Error("/admin/device/import/:id 无效的设备类型ID格式: " + err.Error())
		return
	}

	excelFile, err := c.FormFile("device")
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "获取Excel文件失败: "+err.Error())
		logger.Error("/admin/device/import/:id 获取Excel文件失败: " + err.Error())
		return
	}

	adminID, err := auth.GetCurrentUserID(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "无法获取当前操作的管理员ID"+err.Error())
		logger.Error("/admin/device/import/:id 无法获取当前操作的管理员ID: " + err.Error())
		return
	}

	// 将剩余的工作交由service处理
	err = ctrl.service.UpdateImport(id, adminID, excelFile)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorDeviceNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorDeviceNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/device/import/:id 更新导入设备失败" + err.Error())
		}
		return
	}
	response.Success(c, nil)
}
