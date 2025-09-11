package attachment

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/attachment"
	"xinde/internal/middleware/auth"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

func (ctrl *Controller) FixOrphan(c *gin.Context) {
	var req *dto.FixOrphanReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误: "+err.Error())
		logger.Error("/admin/attachment/fix/orphan 绑定参数错误: " + err.Error())
		return
	}

	// 从上下文中获取当前操作管理员的ID
	userID, err := auth.GetCurrentUserID(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "无法获取当前管理员信息: "+err.Error())
		logger.Error("/admin/attachment/fix/orphan 无法获取当前管理员信息: " + err.Error())
		return
	}
	if !auth.IsAdmin(c) {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "不是管理员没有权限操作!")
		logger.Error("/admin/attachment/fix/orphan 不是管理员没有权限操作!")
		return
	}

	// 将修复任务交由service层处理
	err = ctrl.attachmentService.FixOrphan(userID, req.FilePath, req.Action)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorAttachmentNotFoundOnDesk:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorAttachmentNotFoundOnDesk)
			logger.Error("/admin/attachment/fix/orphan 无法在磁盘上找到孤儿文件")
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/attachment/fix/orphan 修复孤儿文件处理失败: " + err.Error())
		}
		return
	}

	response.Success(c, nil)
}
