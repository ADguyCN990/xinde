package attachment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"xinde/internal/handler/common"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// Delete handles deleting a specific attachment.
// @Summary      删除附件
// @Description  根据附件ID删除对应的文件和数据库记录
// @Tags         Attachment
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "附件 ID"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "删除成功"
// @Failure      400 {object} response.Response "请求参数错误"
// @Failure      401 {object} response.Response "Token错误"
// @Failure      403 {object} response.Response "没有管理员权限"
// @Failure      404 {object} response.Response "附件不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/attachment/{id} [delete]
func (ctrl *Controller) Delete(c *gin.Context) {
	// 从url中获取ID
	id, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorAttachmentIDInvalid)
		logger.Error("DELETE admin/attachment/{id} 无效的ID格式: " + err.Error())
		return
	}
	// 将删除任务交由service层处理
	err = ctrl.attachmentService.Delete(id)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorAttachmentNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorAttachmentNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error(fmt.Sprintf("DELETE admin/attachment/{%d}, 删除附件失败, err: %s", id, err.Error()))
		}
		return
	}
	response.Success(c, nil)
}
