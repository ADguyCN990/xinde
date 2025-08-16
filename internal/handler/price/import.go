package price

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xinde/internal/middleware/auth"
	"xinde/pkg/logger"
	"xinde/pkg/response"
)

// Import handles the import of price data from an Excel file.
// @Summary      导入价格Excel文件
// @Description  上传一个包含价格信息的Excel文件，系统将解析文件内容并批量更新或插入价格数据。如果产品编码已存在，则会用新数据覆盖。
// @Tags         Price
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "要上传的Excel文件 (格式: .xlsx)"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "操作成功"
// @Failure      400 {object} response.Response "文件上传失败或文件内容/格式错误"
// @Failure      401 {object} response.Response "Token错误"
// @Failure      403 {object} response.Response "没有管理员权限"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/prices/import [post]
func (ctrl *Controller) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "文件上传失败: "+err.Error())
		return
	}

	// 从上下文中获取当前操作的管理员ID
	adminID, err := auth.GetCurrentUserID(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "无法获取当前管理员的ID信息: "+err.Error())
		return
	}

	// 调用Service层处理文件
	err = ctrl.priceService.ImportPricesFromFile(c, file, adminID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, err.Error())
		logger.Error("/admin/price/import " + err.Error())
		return
	}

	response.Success(c, nil)
}
