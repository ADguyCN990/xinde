package group

import (
	"github.com/gin-gonic/gin"
	"net/http"
	dto "xinde/internal/dto/group"
	"xinde/internal/handler/common"
	"xinde/internal/middleware/auth"
	"xinde/pkg/logger"
	"xinde/pkg/response"
	"xinde/pkg/stderr"
)

// Update handles the update of an existing group.
// @Summary      更新分组信息
// @Description  根据ID更新一个分组的名称、父级或图标
// @Tags         Group
// @Accept       multipart/form-data
// @Produce      json
// @Param        id   path      int  true  "分组 ID"
// @Param        name formData string false "新的分组名称"
// @Param        parent_id formData int false "新的父级分组ID"
// @Param        icon formData file false "新的分组图标 (可选)"
// @Security     ApiKeyAuth
// @Success      200 {object} response.Response "更新成功"
// @Failure      400 {object} response.Response "请求参数错误或无效ID"
// @Failure      404 {object} response.Response "分组不存在"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /api/v1/admin/group/update/{id} [put]
func (ctrl *Controller) Update(c *gin.Context) {

	groupID, err := common.GetIDFromUrl(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorGroupIDInvalid)
		logger.Error("/admin/group/update 无效的分组ID格式: " + err.Error())
		return
	}

	var req dto.UpdateReq
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, "绑定参数错误")
		logger.Error("/admin/group/update 绑定参数错误: " + err.Error())
		return
	}

	file, _ := c.FormFile("icon")
	adminID, err := auth.GetCurrentUserID(c)
	if err != nil {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "无法获取当前操作的用户ID")
		logger.Error("/admin/group/update 无法获取当前操作的用户ID: " + err.Error())
		return
	}

	// 将剩余的工作交由service处理
	err = ctrl.Service.Update(adminID, groupID, req.ParentID, req.Name, file)
	if err != nil {
		switch err.Error() {
		case stderr.ErrorGroupNotFound:
			response.Error(c, http.StatusNotFound, response.CodeNotFound, "ParentID并不存在")
		case stderr.ErrorCannotMoveGroupIntoItself:
			response.Error(c, http.StatusBadRequest, response.CodeInvalidParams, stderr.ErrorCannotMoveGroupIntoItself)
		default:
			response.Error(c, http.StatusInternalServerError, response.CodeInternalError, stderr.ErrorInternalServerError)
			logger.Error("/admin/group/update 更改分组失败: " + err.Error())
		}
		return
	}
	response.Success(c, nil)
}
