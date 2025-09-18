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

// UpdateImage handles changing the main image of a device type.
// @Summary      修改设备类型图片
// @Description  为一个已存在的设备类型上传新的主图，替换旧图
// @Tags         Device
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   path      int  true  "设备类型 ID"
// @Param        image formData file true "新的设备主图"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "修改成功"
// @Failure      400 {object} response.Response "请求参数错误或无效ID"
// @Failure      404 {object} response.Response "设备类型不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/update/image/{id} [post]
func (ctrl *Controller) UpdateImage(c *gin.Context) {
	deviceTypeID, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorDeviceIDInvalid)
		logger.Error("/admin/device/update/image 无效的设备类型ID格式 " + err.Error())
		return
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "无法获取要更新的图片")
		logger.Error("/admin/device/update/image 无法获取图片文件 " + err.Error())
		return
	}

	adminUID, err := auth.GetCurrentUserID(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "无法获取当前操作的管理员ID")
		logger.Error("/admin/device/update/image 无法获取当前操作的管理员ID: " + err.Error())
		return
	}

	// 剩余的工作交由service处理
	err = ctrl.service.UpdateImage(deviceTypeID, adminUID, imageFile)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorDeviceNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorDeviceNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/device/update/image 更新设备图片发生错误: " + err.Error())
		}
		return
	}
	response.Success(c, nil)
}
