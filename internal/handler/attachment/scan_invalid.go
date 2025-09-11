package attachment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "xinde/internal/dto/attachment"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// ScanInvalid handles scanning for orphan files and records.
// @Summary      扫描异常附件
// @Description  扫描并返回数据库与磁盘文件不一致的记录
// @Tags         Attachment
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200 {object} _.ScanInvalidResp "成功返回扫描结果"
// @Failure      401 {object} response.Response "Token错误"
// @Failure      403 {object} response.Response "没有管理员权限"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/attachment/scan/invalid [get]
func (ctrl *Controller) ScanInvalid(c *gin.Context) {

	orphanData, err := ctrl.attachmentService.ScanInvalid()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
		logger.Error("/admin/attachment/scan/invalid 扫描异常附件失败: " + err.Error())
		return
	}
	response.Success(c, orphanData)
}
