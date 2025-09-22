package device

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/device"
	"xinde/internal/middleware/auth"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// CreateFilterImage handles creating or replacing a filter image configuration.
// @Summary      创建或替换筛选条件图片
// @Description  为一个设备类型下的某个筛选值上传一张图片。如果该配置已存在，则会覆盖旧的图片。
// @Tags         FilterImage
// @Accept       multipart/form-data
// @Produce      json
// @Param        device_type_id formData int true "关联的设备类型ID"
// @Param        filter_value formData string true "匹配的筛选条件的值"
// @Param        image formData file true "要上传的图片文件"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "操作成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "Token错误"
// @Failure      403 {object} response.Response "没有管理员权限"
// @Failure      404 {object} response.Response "关联的设备类型不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/filter_image/create [post]
func (ctrl *Controller) CreateFilterImage(c *gin.Context) {
	var req *dto.CreateFilterImageReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
		logger.Error("/admin/filter_image/create 绑定参数错误" + err.Error())
		return
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "无法获取图片: "+err.Error())
		logger.Error("/admin/filter_image/create 无法获取图片" + err.Error())
		return
	}

	adminUID, err := auth.GetCurrentUserID(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "无法获取管理员ID: "+err.Error())
		logger.Error("/admin/filter_image/create 无法获取管理员ID" + err.Error())
		return
	}

	err = ctrl.service.CreateFilterImage(adminUID, req.DeviceTypeID, req.FilterValue, imageFile)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorDeviceNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorDeviceNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/filter_image/create 创建设备下拉筛选图片失败" + err.Error())
		}
		return
	}

	response.Success(c, nil)
}
