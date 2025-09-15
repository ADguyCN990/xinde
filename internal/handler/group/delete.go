package group

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xinde/internal/handler/common"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// Delete handles the deletion of a group and its descendants.
// @Summary      删除分组
// @Description  删除一个分组及其所有子孙分组，并处理关联的设备
// @Tags         Group
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "分组 ID"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "删除成功"
// @Failure      400 {object} response.Response "请求参数错误或无效ID"
// @Failure      403 {object} response.Response "禁止删除Root分组"
// @Failure      404 {object} response.Response "分组不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/group/delete/{id} [delete]
func (ctrl *Controller) Delete(c *gin.Context) {
	groupID, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorGroupIDInvalid)
		logger.Error("无效的分组ID格式: " + err.Error())
		return
	}

	err = ctrl.Service.Delete(groupID)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorRootGroupCannotBeDeleted:
			response.Error(c, http.StatusForbidden, response.CodeForbidden, stderr.ErrorRootGroupCannotBeDeleted)
		case stderr.ErrorGroupNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, stderr.ErrorGroupNotFound)
		default:
			response.Error(c, http.StatusInternalServerError, http.StatusInternalServerError, stderr.ErrorGroupNotFound)
			logger.Error("/admin/group/delete 删除分组发生错误: " + err.Error())
		}
		return
	}
	response.Success(c, nil)
}
