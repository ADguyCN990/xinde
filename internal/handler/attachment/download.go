package attachment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"xinde/internal/handler/common"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
	"xinde/pkg/util"
)

// Download handles downloading a specific attachment.
// @Summary      下载附件
// @Description  根据附件ID下载对应的文件
// @Tags         Attachment
// @Produce      application/octet-stream
// @Param        id   path      int  true  "附件 ID"
// @Security     ApiKeyAuth
// @Success      200 {file} file "文件流"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "Token错误"
// @Failure      403 {object} response.Response "没有管理员权限"
// @Failure      404 {object} response.Response "附件或文件不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/attachments/{id}/download [get]
func (ctrl *Controller) Download(c *gin.Context) {
	// 获取ID
	id, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorAttachmentIDInvalid)
		logger.Error("/admin/attachment/download 无效的ID格式: " + err.Error())
		return
	}

	fileName, fileType, file, err := ctrl.attachmentService.GetAttachmentForDownload(id)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorAttachmentNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorAttachmentNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "服务器内部错误")
			logger.Error("/admin/attachment/download 获取附件失败: " + err.Error())
		}
		return
	}
	defer file.Close()

	// 为了下载，设置http响应头
	c.Header("Content-Disposition", util.FormatContentDisposition(fileName))
	c.Header("Content-Type", fileType)

	// 将文件流写入响应体
	_, err = io.Copy(c.Writer, file)
	if err != nil {
		logger.Error(fmt.Sprintf("/admin/attachment/download 向客户端写入文件流时出错! attachmentID: %d, err: %s", id, err.Error()))
	}
}
