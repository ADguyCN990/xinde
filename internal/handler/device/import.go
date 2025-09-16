package device

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/device"
	"xinde/internal/middleware/auth"
	"xinde/internal/service/device"
	"xinde/pkg/logger"
	"xinde/pkg/response"
)

type Controller struct {
	service *device.Service
}

func NewDeviceController() (*Controller, error) {
	service, err := device.NewDeviceService()
	if err != nil {
		return nil, fmt.Errorf("创建service实例失败: " + err.Error())
	}
	return &Controller{service: service}, nil
}

// Import handles importing devices from an Excel file.
// @Summary      从Excel导入设备
// @Description  上传Excel文件，批量导入设备（方案）到一个指定分组
// @Tags         Device
// @Accept       multipart/form-data
// @Produce      json
// @Param        group_id formData int true "目标分组ID"
// @Param        excel formData file true "包含设备数据的Excel文件"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "导入成功，并返回导入数量"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      500 {object} response.Response "服务器内部错误或导入失败"
// @Router       /api/v1/admin/device/import [post]
func (ctrl *Controller) Import(c *gin.Context) {
	var req dto.ImportReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数失败: "+err.Error())
		logger.Error("/admin/device/import 绑定参数失败: " + err.Error())
		return
	}

	excelFile, err := c.FormFile("device")
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "获取上传的Excel文件失败: "+err.Error())
		logger.Error("/admin/device/import 获取上传的excel文件失败: " + err.Error())
		return
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "获取上传的主图失败: "+err.Error())
		logger.Error("/admin/device/import 获取上传的主图失败: " + err.Error())
		return
	}

	adminID, err := auth.GetCurrentUserID(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "无法获取当前操作的管理员ID"+err.Error())
		logger.Error("/admin/device/import 无法获取当前操作的管理员ID: " + err.Error())
		return
	}

	// 将剩余的工作交由service处理
	ctrl.service.ImportFromExcel(adminID, req.GroupID, req.DeviceTypeName, excelFile, imageFile)
}
